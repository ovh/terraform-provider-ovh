package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudStorageObjectBucketObjectLockResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudStorageObjectBucketObjectLockResource)(nil)
	_ resource.ResourceWithImportState = (*cloudStorageObjectBucketObjectLockResource)(nil)
)

func NewCloudStorageObjectBucketObjectLockResource() resource.Resource {
	return &cloudStorageObjectBucketObjectLockResource{}
}

type cloudStorageObjectBucketObjectLockResource struct {
	config *Config
}

// cloudStorageObjectBucketObjectLockModel is the Terraform state model for the
// object-lock resource.
type cloudStorageObjectBucketObjectLockModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	BucketName  ovhtypes.TfStringValue `tfsdk:"bucket_name"`
	ObjectKey   ovhtypes.TfStringValue `tfsdk:"object_key"`
	VersionId   ovhtypes.TfStringValue `tfsdk:"version_id"`

	Retention types.Object `tfsdk:"retention"`
	LegalHold types.Object `tfsdk:"legal_hold"`

	Id ovhtypes.TfStringValue `tfsdk:"id"`
}

func (r *cloudStorageObjectBucketObjectLockResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket_object_lock"
}

func (r *cloudStorageObjectBucketObjectLockResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.config = config
}

func (r *cloudStorageObjectBucketObjectLockResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages retention (object lock) and legal hold settings on a single S3 object " +
			"(optionally a specific version) using the v2 API. " +
			"Destroying this resource removes it from Terraform state but does NOT remove the retention " +
			"or legal hold from the object itself: S3 object lock, once set to COMPLIANCE or GOVERNANCE, " +
			"cannot be cleared through the public cloud API.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"bucket_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Name of the bucket",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"object_key": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Object key (path within the bucket)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Optional:    true,
				Description: "Version identifier of the object. When set, the resource manages the lock on that specific version.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"retention": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Retention (object lock) configuration to apply. Only COMPLIANCE and GOVERNANCE modes can be set; NONE is rejected by the API.",
				Attributes: map[string]schema.Attribute{
					"mode": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Required:    true,
						Description: "Retention mode (COMPLIANCE or GOVERNANCE).",
					},
					"retain_until_date": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Required:    true,
						Description: "RFC 3339 date/time until which the object is retained.",
					},
				},
			},
			"legal_hold": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Legal hold configuration to apply.",
				Attributes: map[string]schema.Attribute{
					"status": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Required:    true,
						Description: "Legal hold status (ON or OFF).",
					},
				},
			},

			"id": schema.StringAttribute{
				CustomType: ovhtypes.TfStringType{},
				Computed:   true,
				Description: "Resource identifier formatted as " +
					"`<service_name>/<bucket_name>/<url-encoded object_key>` " +
					"or `<service_name>/<bucket_name>/<url-encoded object_key>/<version_id>` " +
					"when a version is targeted.",
			},
		},
	}
}

// composeObjectLockID builds the composite ID for the resource. The object key
// is URL-encoded so that slashes inside the key cannot be mistaken for
// separators.
func composeObjectLockID(serviceName, bucketName, objectKey, versionId string) string {
	id := serviceName + "/" + bucketName + "/" + url.PathEscape(objectKey)
	if versionId != "" {
		id += "/" + versionId
	}
	return id
}

func (r *cloudStorageObjectBucketObjectLockResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 3 && len(splits) != 4 {
		resp.Diagnostics.AddError(
			"Given ID is malformed",
			"ID must be formatted like: <service_name>/<bucket_name>/<url-encoded object_key>[/<version_id>]",
		)
		return
	}

	serviceName := splits[0]
	bucketName := splits[1]
	encodedKey := splits[2]
	objectKey, err := url.PathUnescape(encodedKey)
	if err != nil {
		resp.Diagnostics.AddError("Given ID is malformed", fmt.Sprintf("unable to URL-decode object key %q: %s", encodedKey, err))
		return
	}

	var versionId string
	if len(splits) == 4 {
		versionId = splits[3]
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), serviceName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("bucket_name"), bucketName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("object_key"), objectKey)...)
	if versionId != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("version_id"), versionId)...)
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// objectEndpoint builds the apiv2 endpoint for the object (or version) managed
// by this resource.
func objectEndpoint(serviceName, bucketName, objectKey, versionId string) string {
	base := "/v2/publicCloud/project/" + url.PathEscape(serviceName) +
		"/storage/object/bucket/" + url.PathEscape(bucketName) +
		"/object/" + url.PathEscape(objectKey)
	if versionId != "" {
		base += "/version/" + url.PathEscape(versionId)
	}
	return base
}

// buildUpdatePayload converts the plan into the API payload.
func buildUpdatePayload(data cloudStorageObjectBucketObjectLockModel) *CloudBucketObjectUpdatePayload {
	payload := &CloudBucketObjectUpdatePayload{}

	if !data.Retention.IsNull() && !data.Retention.IsUnknown() {
		ra := data.Retention.Attributes()
		ret := &CloudBucketObjectRetention{}
		if v, ok := ra["mode"].(ovhtypes.TfStringValue); ok {
			ret.Mode = v.ValueString()
		}
		if v, ok := ra["retain_until_date"].(ovhtypes.TfStringValue); ok {
			ret.RetainUntil = v.ValueString()
		}
		payload.Retention = ret
	}

	if !data.LegalHold.IsNull() && !data.LegalHold.IsUnknown() {
		la := data.LegalHold.Attributes()
		if v, ok := la["status"].(ovhtypes.TfStringValue); ok {
			payload.LegalHold = v.ValueString()
		}
	}

	return payload
}

// mergeObjectResponse populates retention and legal_hold from the API response.
// Unlike the simple pass-through in other resources, we preserve the user-
// provided configuration when the API did not return an explicit value (an
// object created before object-lock was set up on the bucket, for example,
// reports no retention or legal hold — but the user may still have asked for
// one, and removing it would create a plan loop).
func (m *cloudStorageObjectBucketObjectLockModel) mergeObjectResponse(o *CloudBucketObjectAPI) {
	if o == nil {
		m.Retention = types.ObjectNull(cloudBucketObjectRetentionAttrTypes())
		m.LegalHold = types.ObjectNull(cloudBucketObjectLegalHoldAttrTypes())
		return
	}

	if o.Retention != nil {
		m.Retention = apiObjectRetentionToTFObject(o.Retention)
	} else {
		m.Retention = types.ObjectNull(cloudBucketObjectRetentionAttrTypes())
	}

	if o.LegalHold != "" {
		m.LegalHold = apiLegalHoldStringToTFObject(o.LegalHold)
	} else {
		m.LegalHold = types.ObjectNull(cloudBucketObjectLegalHoldAttrTypes())
	}
}

func (r *cloudStorageObjectBucketObjectLockResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data cloudStorageObjectBucketObjectLockModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := objectEndpoint(
		data.ServiceName.ValueString(),
		data.BucketName.ValueString(),
		data.ObjectKey.ValueString(),
		data.VersionId.ValueString(),
	)

	payload := buildUpdatePayload(data)

	var response CloudBucketObjectAPI
	if err := r.config.OVHClient.Put(endpoint, payload, &response); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	data.mergeObjectResponse(&response)
	data.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(composeObjectLockID(
		data.ServiceName.ValueString(),
		data.BucketName.ValueString(),
		data.ObjectKey.ValueString(),
		data.VersionId.ValueString(),
	))}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketObjectLockResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data cloudStorageObjectBucketObjectLockModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := objectEndpoint(
		data.ServiceName.ValueString(),
		data.BucketName.ValueString(),
		data.ObjectKey.ValueString(),
		data.VersionId.ValueString(),
	)

	var response CloudBucketObjectAPI
	if err := r.config.OVHClient.Get(endpoint, &response); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.mergeObjectResponse(&response)
	data.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(composeObjectLockID(
		data.ServiceName.ValueString(),
		data.BucketName.ValueString(),
		data.ObjectKey.ValueString(),
		data.VersionId.ValueString(),
	))}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketObjectLockResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state cloudStorageObjectBucketObjectLockModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := objectEndpoint(
		plan.ServiceName.ValueString(),
		plan.BucketName.ValueString(),
		plan.ObjectKey.ValueString(),
		plan.VersionId.ValueString(),
	)

	payload := buildUpdatePayload(plan)

	var response CloudBucketObjectAPI
	if err := r.config.OVHClient.Put(endpoint, payload, &response); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	plan.mergeObjectResponse(&response)
	plan.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(composeObjectLockID(
		plan.ServiceName.ValueString(),
		plan.BucketName.ValueString(),
		plan.ObjectKey.ValueString(),
		plan.VersionId.ValueString(),
	))}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete intentionally removes the resource only from Terraform state. S3 object
// lock is immutable once set (COMPLIANCE cannot be cleared at all before the
// retain-until date; GOVERNANCE would require a bypass header that the apiv2
// does not expose). Legal hold could in principle be cleared via PUT, but
// removing the resource does not necessarily express the intent to disable it,
// so we keep the backend untouched.
func (r *cloudStorageObjectBucketObjectLockResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// no-op: see method comment
}

package ovh

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudStorageObjectBucketACLResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudStorageObjectBucketACLResource)(nil)
	_ resource.ResourceWithImportState = (*cloudStorageObjectBucketACLResource)(nil)
)

func NewCloudStorageObjectBucketACLResource() resource.Resource {
	return &cloudStorageObjectBucketACLResource{}
}

type cloudStorageObjectBucketACLResource struct {
	config *Config
}

func (r *cloudStorageObjectBucketACLResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket_acl"
}

func (r *cloudStorageObjectBucketACLResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func storageObjectBucketCannedACLs() []string {
	names := make([]string, 0, len(s3CannedACLs))
	for name := range s3CannedACLs {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (r *cloudStorageObjectBucketACLResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	cannedACLs := strings.Join(storageObjectBucketCannedACLs(), ", ")

	resp.Schema = schema.Schema{
		Description: "Manages the canned ACL of an S3 compatible object storage bucket in a public cloud project, through the S3 API " +
			"with the S3 credentials of the provider configuration. Destroying it sets the bucket back to private.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "Service name of the resource representing the id of the cloud project. If omitted, the OVH_CLOUD_PROJECT_SERVICE environment variable is used.",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
				PlanModifiers: []planmodifier.String{
					EnvDefaultString("OVH_CLOUD_PROJECT_SERVICE", true),
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"bucket_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Bucket identifier, in the <REGION>_<name> form",
				MarkdownDescription: "Bucket identifier, in the `<REGION>_<name>` form",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "ACL identifier, the bucket identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"acl": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Canned ACL of the bucket (" + cannedACLs + ")",
				Validators: []validator.String{
					stringvalidator.OneOf(storageObjectBucketCannedACLs()...),
				},
			},
		},
	}
}

func (r *cloudStorageObjectBucketACLResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<bucket_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("bucket_id"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

// resolve returns the S3 client and where the bucket is served. The error of a missing bucket is the 404 of the publicCloud API.
func (r *cloudStorageObjectBucketACLResource) resolve(ctx context.Context, data *CloudStorageObjectBucketACLModel) (*s3Client, *storageObjectBucketS3Target, error) {
	client, err := r.config.S3Client()
	if err != nil {
		return nil, nil, err
	}

	target, err := storageObjectBucketS3TargetRead(ctx, r.config.OVHClient, data.ServiceName.ValueString(), data.BucketId.ValueString())
	if err != nil {
		return nil, nil, err
	}

	return client, target, nil
}

func (r *cloudStorageObjectBucketACLResource) apply(ctx context.Context, data *CloudStorageObjectBucketACLModel) diag.Diagnostics {
	var diags diag.Diagnostics

	client, target, err := r.resolve(ctx, data)
	if err != nil {
		diags.AddError(fmt.Sprintf("Cannot reach bucket %s on the S3 API", data.BucketId.ValueString()), err.Error())
		return diags
	}

	if err := client.PutBucketCannedACL(ctx, target.Host, target.Region, data.ACL.ValueString()); err != nil {
		diags.AddError(fmt.Sprintf("Error setting the ACL of bucket %s", target.Host), err.Error())
		return diags
	}

	data.Id = data.BucketId
	return diags
}

func (r *cloudStorageObjectBucketACLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudStorageObjectBucketACLModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.apply(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketACLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudStorageObjectBucketACLModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, target, err := r.resolve(ctx, &data)
	if err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		endpoint := storageObjectBucketEndpoint(data.ServiceName.ValueString(), data.BucketId.ValueString())
		if storageObjectBucketReadUnavailable(err, endpoint, &resp.Diagnostics) {
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot reach bucket %s on the S3 API", data.BucketId.ValueString()), err.Error())
		return
	}

	policy, err := client.GetBucketACL(ctx, target.Host, target.Region)
	if err != nil {
		if isS3NoSuchBucket(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error reading the ACL of bucket %s", target.Host), err.Error())
		return
	}

	// Grants that amount to no canned ACL cannot be expressed here: the state is kept, the next apply sets the canned ACL back.
	if acl := policy.CannedACL(); acl != "" {
		data.ACL = ovhtypes.NewTfStringValue(acl)
	} else {
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("The ACL of bucket %s matches no canned ACL", target.Host),
			fmt.Sprintf("The bucket has %d grants that amount to no canned ACL. The state keeps acl = %q, the next apply sets it back.", len(policy.Grants), data.ACL.ValueString()),
		)
	}
	data.Id = data.BucketId

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketACLResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CloudStorageObjectBucketACLModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.apply(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketACLResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudStorageObjectBucketACLModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, target, err := r.resolve(ctx, &data)
	if err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot reach bucket %s on the S3 API", data.BucketId.ValueString()), err.Error())
		return
	}

	if err := client.PutBucketCannedACL(ctx, target.Host, target.Region, "private"); err != nil && !isS3NoSuchBucket(err) {
		resp.Diagnostics.AddError(fmt.Sprintf("Error setting the ACL of bucket %s back to private", target.Host), err.Error())
	}
}

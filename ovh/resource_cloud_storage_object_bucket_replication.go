package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudStorageObjectBucketReplicationResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudStorageObjectBucketReplicationResource)(nil)
	_ resource.ResourceWithImportState = (*cloudStorageObjectBucketReplicationResource)(nil)
)

func NewCloudStorageObjectBucketReplicationResource() resource.Resource {
	return &cloudStorageObjectBucketReplicationResource{}
}

type cloudStorageObjectBucketReplicationResource struct {
	config *Config
}

func (r *cloudStorageObjectBucketReplicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket_replication"
}

func (r *cloudStorageObjectBucketReplicationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cloudStorageObjectBucketReplicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the replication configuration of an S3 object storage bucket using the v2 API. " +
			"Setting an empty rules list removes all replication rules.",
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
				Description: "Name of the source bucket",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Resource identifier (service_name/bucket_name)",
			},
			"rules": schema.ListNestedAttribute{
				Required:    true,
				Description: "List of replication rules. An empty list removes all rules.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Required:    true,
							Description: "Unique identifier for the rule",
						},
						"status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Required:    true,
							Description: "Rule status (ENABLED or DISABLED)",
						},
						"priority": schema.Int64Attribute{
							Required:    true,
							Description: "Rule priority (higher value = higher priority)",
						},
						"filter": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "Filter that identifies objects to replicate",
							Attributes: map[string]schema.Attribute{
								"prefix": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Optional:    true,
									Description: "Object key prefix",
								},
								"tags": schema.MapAttribute{
									ElementType: ovhtypes.TfStringType{},
									Optional:    true,
									Description: "Object tags to match",
								},
							},
						},
						"destination": schema.SingleNestedAttribute{
							Required:    true,
							Description: "Destination bucket for replicated objects",
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Required:    true,
									Description: "Destination bucket name",
								},
								"region": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Required:    true,
									Description: "Destination bucket region",
								},
								"storage_class": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Optional:    true,
									Description: "Storage class for replicated objects",
								},
							},
						},
						"delete_marker_replication": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Required:    true,
							Description: "Whether to replicate delete markers (ENABLED or DISABLED)",
						},
					},
				},
			},
		},
	}
}

func (r *cloudStorageObjectBucketReplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like: <service_name>/<bucket_name>")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("bucket_name"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func (r *cloudStorageObjectBucketReplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudStorageObjectBucketReplicationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.putReplication(data); err != nil {
		resp.Diagnostics.AddError("Error setting replication rules", err.Error())
		return
	}

	if err := r.readReplication(&data); err != nil {
		resp.Diagnostics.AddError("Error reading replication rules", err.Error())
		return
	}

	data.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(data.ServiceName.ValueString() + "/" + data.BucketName.ValueString())}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketReplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudStorageObjectBucketReplicationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readReplication(&data); err != nil {
		resp.Diagnostics.AddError("Error reading replication rules", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketReplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state CloudStorageObjectBucketReplicationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Id = state.Id

	if err := r.putReplication(data); err != nil {
		resp.Diagnostics.AddError("Error updating replication rules", err.Error())
		return
	}

	if err := r.readReplication(&data); err != nil {
		resp.Diagnostics.AddError("Error reading replication rules after update", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketReplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudStorageObjectBucketReplicationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete all rules by sending an empty list
	emptyRules, _ := types.ListValue(types.ObjectType{AttrTypes: replicationRuleAttrTypes()}, []attr.Value{})
	data.Rules = emptyRules

	if err := r.putReplication(data); err != nil {
		resp.Diagnostics.AddError("Error deleting replication rules", err.Error())
	}
}

func (r *cloudStorageObjectBucketReplicationResource) putReplication(data CloudStorageObjectBucketReplicationModel) error {
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/storage/object/bucket/" + url.PathEscape(data.BucketName.ValueString()) +
		"/replication"

	payload := CloudBucketReplicationRequest{
		Rules: tfListToAPIReplicationRules(data.Rules),
	}

	var response CloudBucketReplicationResponse
	if err := r.config.OVHClient.Put(endpoint, payload, &response); err != nil {
		return fmt.Errorf("calling Put %s: %w", endpoint, err)
	}
	return nil
}

func (r *cloudStorageObjectBucketReplicationResource) readReplication(data *CloudStorageObjectBucketReplicationModel) error {
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/storage/object/bucket/" + url.PathEscape(data.BucketName.ValueString()) +
		"/replication"

	var response CloudBucketReplicationResponse
	if err := r.config.OVHClient.Get(endpoint, &response); err != nil {
		return fmt.Errorf("calling Get %s: %w", endpoint, err)
	}

	data.Rules = apiReplicationRulesToTFList(response.Rules, data.Rules)
	return nil
}

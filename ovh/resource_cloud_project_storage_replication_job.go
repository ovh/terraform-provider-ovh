package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ resource.ResourceWithConfigure = (*cloudProjectStorageReplicationJobResource)(nil)

func NewCloudProjectStorageReplicationJobResource() resource.Resource {
	return &cloudProjectStorageReplicationJobResource{}
}

type cloudProjectStorageReplicationJobResource struct {
	config *Config
}

func (r *cloudProjectStorageReplicationJobResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_storage_replication_job"
}

func (r *cloudProjectStorageReplicationJobResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cloudProjectStorageReplicationJobResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CloudProjectStorageReplicationJobResourceSchema(ctx)
}

func (r *cloudProjectStorageReplicationJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudProjectStorageReplicationJobModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Trigger the replication job
	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/storage/%s/job/replication",
		url.PathEscape(data.ServiceName.ValueString()),
		url.PathEscape(data.RegionName.ValueString()),
		url.PathEscape(data.ContainerName.ValueString()),
	)

	var responseData CloudProjectStorageReplicationJobResponse
	if err := r.config.OVHClient.Post(endpoint, nil, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Set the ID from the API response
	data.ID = ovhtypes.NewTfStringValue(responseData.ID)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudProjectStorageReplicationJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// No-op: The API doesn't support retrieving job status
	// We just keep whatever is in state
}

func (r *cloudProjectStorageReplicationJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No-op: All attributes have RequiresReplace, so Update should never be called
	resp.Diagnostics.AddError(
		"Update not supported",
		"All attributes require replacement. This should not happen.",
	)
}

func (r *cloudProjectStorageReplicationJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op: Replication jobs cannot be cancelled or deleted once triggered
	// Simply remove from state
}

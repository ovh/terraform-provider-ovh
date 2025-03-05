package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

var _ resource.ResourceWithConfigure = (*cloudProjectVolumeBackupResource)(nil)

func NewCloudProjectVolumeBackupResource() resource.Resource {
	return &cloudProjectVolumeBackupResource{}
}

type cloudProjectVolumeBackupResource struct {
	config *Config
}

func (r *cloudProjectVolumeBackupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_volume_backup"
}

func (d *cloudProjectVolumeBackupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	d.config = config
}

func (d *cloudProjectVolumeBackupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CloudProjectVolumeBackupResourceSchema(ctx)
}

func (r *cloudProjectVolumeBackupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 3 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<region_name>/<id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("region_name"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[2])...)
}

func (r *cloudProjectVolumeBackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData CloudProjectVolumeBackupModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create volume backup
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.RegionName.ValueString()) + "/volumeBackup"
	if err := r.config.OVHClient.PostWithContext(ctx, endpoint, data.ToCreate(), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for backup to be in status "ok"
	endpoint += "/" + url.PathEscape(responseData.Id.ValueString())
	err := retry.RetryContext(ctx, 30*time.Minute, func() *retry.RetryError {
		if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
			return retry.NonRetryableError(fmt.Errorf("call to %q failed, retrying (error: %s)", endpoint, err))
		}

		switch responseData.Status.ValueString() {
		case "error":
			return retry.NonRetryableError(errors.New("resource is in status ERROR"))
		case "ok":
			return nil
		default:
			return retry.RetryableError(fmt.Errorf("resource not ready, retrying (current status: %s)", responseData.Status.ValueString()))
		}
	})
	if err != nil {
		resp.Diagnostics.AddError("error waiting for backup to be created", err.Error())
		return
	}

	responseData.MergeWith(&data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectVolumeBackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData CloudProjectVolumeBackupModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.RegionName.ValueString()) +
		"/volumeBackup/" + url.PathEscape(data.Id.ValueString())

	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudProjectVolumeBackupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("not implemented", "this function should never be called")
}

func (r *cloudProjectVolumeBackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudProjectVolumeBackupModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.RegionName.ValueString()) +
		"/volumeBackup/" + url.PathEscape(data.Id.ValueString())
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
	}
}

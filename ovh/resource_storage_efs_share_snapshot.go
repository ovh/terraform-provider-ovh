package ovh

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

var _ resource.ResourceWithConfigure = (*storageEfsShareSnapshotResource)(nil)

func NewStorageEfsShareSnapshotResource() resource.Resource {
	return &storageEfsShareSnapshotResource{}
}

type storageEfsShareSnapshotResource struct {
	config *Config
}

func (r *storageEfsShareSnapshotResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_efs_share_snapshot"
}

func (d *storageEfsShareSnapshotResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *storageEfsShareSnapshotResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = StorageEfsShareSnapshotResourceSchema(ctx)
}

func (r *storageEfsShareSnapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData StorageEfsShareSnapshotModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share/" + url.PathEscape(data.ShareId.ValueString()) + "/snapshot"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	resWait, err := r.WaitForSnapshotCreation(ctx, r.config.OVHClient, data.ServiceName.ValueString(), data.ShareId.ValueString(), responseData.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Operation"),
			err.Error(),
		)
		return

	}

	resSnapshot := &StorageEfsShareSnapshotModel{}
	endpoint = "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share/" + url.PathEscape(data.ShareId.ValueString()) + "/snapshot/" + url.PathEscape(resWait.(string))
	if err := r.config.OVHClient.Get(endpoint, resSnapshot); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(resSnapshot)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *storageEfsShareSnapshotResource) WaitForSnapshotCreation(ctx context.Context, client *ovhwrap.Client, serviceName, shareID, snapshotID string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"creating"},
		Target:  []string{"available"},
		Refresh: func() (interface{}, string, error) {
			res := &StorageEfsShareSnapshotModel{}
			endpoint := "/storage/netapp/" + url.PathEscape(serviceName) + "/share/" + url.PathEscape(shareID) + "/snapshot/" + url.PathEscape(snapshotID)
			err := client.GetWithContext(ctx, endpoint, res)
			if err != nil {
				return res, "", err
			}
			return res.Id.ValueString(), res.Status.ValueString(), nil
		},
		Timeout:    360 * time.Second,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	res, err := stateConf.WaitForStateContext(ctx)

	return res, err
}

func (r *storageEfsShareSnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData StorageEfsShareSnapshotModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share/" + url.PathEscape(data.ShareId.ValueString()) + "/snapshot/" + url.PathEscape(data.Id.ValueString()) + ""

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

func (r *storageEfsShareSnapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData, responseData StorageEfsShareSnapshotModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update resource
	endpoint := "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share/" + url.PathEscape(data.ShareId.ValueString()) + "/snapshot/" + url.PathEscape(data.Id.ValueString()) + ""
	if err := r.config.OVHClient.Put(endpoint, planData.ToUpdate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Read updated resource
	endpoint = "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share/" + url.PathEscape(data.ShareId.ValueString()) + "/snapshot/" + url.PathEscape(data.Id.ValueString()) + ""
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&planData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *storageEfsShareSnapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data StorageEfsShareSnapshotModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share/" + url.PathEscape(data.ShareId.ValueString()) + "/snapshot/" + url.PathEscape(data.Id.ValueString()) + ""
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
	}
}

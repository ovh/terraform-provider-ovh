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

var _ resource.ResourceWithConfigure = (*storageEfsShareResource)(nil)

func NewStorageEfsShareResource() resource.Resource {
	return &storageEfsShareResource{}
}

type storageEfsShareResource struct {
	config *Config
}

func (r *storageEfsShareResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_efs_share"
}

func (d *storageEfsShareResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *storageEfsShareResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = StorageEfsShareResourceSchema(ctx)
}

func (r *storageEfsShareResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData StorageEfsShareModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	resWait, err := r.WaitForShareCreation(ctx, r.config.OVHClient, data.ServiceName.ValueString(), responseData.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error calling Operation",
			err.Error(),
		)
		return
	}

	resShare := &StorageEfsShareModel{}
	endpoint = "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share/" + url.PathEscape(resWait.(string))
	if err := r.config.OVHClient.Get(endpoint, resShare); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(resShare)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *storageEfsShareResource) WaitForShareCreation(ctx context.Context, client *ovhwrap.Client, serviceName, shareID string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"creating", "creating_from_snapshot"},
		Target:  []string{"available"},
		Refresh: func() (interface{}, string, error) {
			res := &StorageEfsShareModel{}
			endpoint := "/storage/netapp/" + url.PathEscape(serviceName) + "/share/" + url.PathEscape(shareID)
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

func (r *storageEfsShareResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData StorageEfsShareModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share/" + url.PathEscape(data.Id.ValueString()) + ""

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

func (r *storageEfsShareResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData, responseData StorageEfsShareModel

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
	endpoint := "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share/" + url.PathEscape(data.Id.ValueString())
	if err := r.config.OVHClient.Put(endpoint, planData.ToUpdate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Check if size has been modified as updating it requires a specific call
	if !planData.Size.Equal(data.Size) {
		endpoint := "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share/" + url.PathEscape(data.Id.ValueString())
		suffix := "/extend"
		if planData.Size.ValueInt64() < data.Size.ValueInt64() {
			suffix = "/shrink"
		}
		endpoint = endpoint + suffix

		if err := r.config.OVHClient.Post(endpoint, map[string]any{
			"size": planData.Size.ValueInt64(),
		}, nil); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Post %s", endpoint),
				err.Error(),
			)
			return
		}

		_, err := r.WaitForShareSizeUpdate(ctx, r.config.OVHClient, data.ServiceName.ValueString(), data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Post %s", endpoint),
				err.Error(),
			)
			return
		}
	}

	// Read updated resource
	endpoint = "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share/" + url.PathEscape(data.Id.ValueString())
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&planData)
	responseData.MergeWith(&data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *storageEfsShareResource) WaitForShareSizeUpdate(ctx context.Context, client *ovhwrap.Client, serviceName, shareID string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"shrinking", "extending"},
		Target:  []string{"available"},
		Refresh: func() (interface{}, string, error) {
			res := &StorageEfsModel{}
			endpoint := "/storage/netapp/" + url.PathEscape(serviceName) + "/share/" + url.PathEscape(shareID)
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

func (r *storageEfsShareResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data StorageEfsShareModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share/" + url.PathEscape(data.Id.ValueString()) + ""
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
	}
}

package ovh

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
)

var _ resource.ResourceWithConfigure = (*cloudProjectVolumeResource)(nil)

func NewCloudProjectVolumeResource() resource.Resource {
	return &cloudProjectVolumeResource{}
}

type cloudProjectVolumeResource struct {
	config *Config
}

func (r *cloudProjectVolumeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_volume"
}

func (d *cloudProjectVolumeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *cloudProjectVolumeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CloudProjectVolumeResourceSchema(ctx)
}

func (r *cloudProjectVolumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData CloudProjectVolumeModelOp

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/volume"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	resWait, err := r.WaitForVolumeCreation(ctx, r.config.OVHClient, data.ServiceName.ValueString(), responseData.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Operation"),
			err.Error(),
		)
		return
	}

	resVol := &CloudProjectVolumeModelOp{}
	endpoint = "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/volume/" + url.PathEscape(resWait.(string))
	if err := r.config.OVHClient.Get(endpoint, resVol); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(resVol)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *cloudProjectVolumeResource) WaitForVolumeCreation(ctx context.Context, client *ovh.Client, serviceName, operationId string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"null", "in-progress", "created", ""},
		Target:  []string{"completed"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectVolumeModelOp{}
			endpoint := "/cloud/project/" + url.PathEscape(serviceName) + "/operation/" + url.PathEscape(operationId)
			err := client.GetWithContext(ctx, endpoint, res)
			if err != nil {
				return res, "", err
			}
			return res.ResourceId.ValueString(), res.Status.ValueString(), nil
		},
		Timeout:    360 * time.Second,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	res, err := stateConf.WaitForStateContext(ctx)

	return res, err
}

func (r *cloudProjectVolumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData CloudProjectVolumeModelOp

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/volume/" + url.PathEscape(data.VolumeId.ValueString())

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

func (r *cloudProjectVolumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData, responseData CloudProjectVolumeModelOp

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
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/volume/" + url.PathEscape(data.VolumeId.ValueString())
	if err := r.config.OVHClient.Put(endpoint, planData.ToUpdate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Check if size has been modified as updating it requires a specific call
	if !planData.Size.Equal(data.Size) {
		endpoint = "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/volume/" + url.PathEscape(data.VolumeId.ValueString()) + "/upsize"
		if err := r.config.OVHClient.Post(endpoint, map[string]any{
			"size": planData.Size.ValueInt64(),
		}, nil); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Post %s", endpoint),
				err.Error(),
			)
			return
		}

		// Here we need to wait for some time to make sure the size is correctly updated on backend side
		time.Sleep(1 * time.Minute)
	}

	// Read updated resource
	endpoint = "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/volume/" + url.PathEscape(data.VolumeId.ValueString())
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

func (r *cloudProjectVolumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudProjectVolumeModelOp

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/volume/" + url.PathEscape(data.VolumeId.ValueString())
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
	}
}

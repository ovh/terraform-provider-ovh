package ovh

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ resource.ResourceWithConfigure = (*cloudProjectLoadbalancerResource)(nil)

func NewCloudProjectLoadbalancerResource() resource.Resource {
	return &cloudProjectLoadbalancerResource{}
}

type cloudProjectLoadbalancerResource struct {
	config *Config
}

func (r *cloudProjectLoadbalancerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_loadbalancer"
}

func (d *cloudProjectLoadbalancerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cloudProjectLoadbalancerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 3 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<region_name>/<id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("region_name"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[2])...)

	// Setting a placeholder for required field `network` because
	// it is not returned by the API at read time
	networkPlaceholder := NetworkValue{
		Private: NetworkPrivateValue{
			Network: NetworkPrivateNetworkValue{
				Id:       types.NewTfStringValue("id_placeholder"),
				SubnetId: types.NewTfStringValue("subnet_id_placeholder"),
				state:    attr.ValueStateKnown,
			},
			state: attr.ValueStateKnown,
		},
		state: attr.ValueStateKnown,
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network"), networkPlaceholder)...)
}

func (d *cloudProjectLoadbalancerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CloudProjectLoadbalancerResourceSchema(ctx)
}

func (r *cloudProjectLoadbalancerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var (
		data, responseData CloudProjectRegionLoadbalancerModel
		operation          CloudProjectOperationResponse
	)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle service_name: use provided value or fall back to environment variable
	if data.ServiceName.IsNull() || data.ServiceName.IsUnknown() || data.ServiceName.ValueString() == "" {
		envServiceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE")
		if envServiceName == "" {
			resp.Diagnostics.AddError(
				"Missing service_name",
				"The service_name attribute is required. Please provide it in the resource configuration or set the OVH_CLOUD_PROJECT_SERVICE environment variable.",
			)
			return
		}
		data.ServiceName = types.NewTfStringValue(envServiceName)
	}

	// Create resource and get operation ID
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/loadbalancing/loadbalancer"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &operation); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for operation to be completed
	resourceID, err := waitForCloudProjectOperation(ctx, r.config.OVHClient, data.ServiceName.ValueString(), operation.Id, "", defaultCloudOperationTimeout)
	if err != nil {
		resp.Diagnostics.AddError("error waiting for operation", err.Error())
	}

	// Update resource
	endpoint = "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/loadbalancing/loadbalancer/" + url.PathEscape(resourceID)
	if err := r.config.OVHClient.Put(endpoint, data.ToUpdate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Read updated resource
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectLoadbalancerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData CloudProjectRegionLoadbalancerModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/loadbalancing/loadbalancer/" + url.PathEscape(data.Id.ValueString())

	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Explicitely set the 'listeners' attribute from the state as it is not returned by the API
	savedListeners := data.Listeners
	data.MergeWith(&responseData)
	data.Listeners = savedListeners

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudProjectLoadbalancerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudProjectRegionLoadbalancerModel

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
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/loadbalancing/loadbalancer/" + url.PathEscape(data.Id.ValueString())
	if err := r.config.OVHClient.Put(endpoint, planData.ToUpdate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for flavor to be updated (if changed) and read updated data
	responseData, err := r.waitForLoadBalancerToBeReady(ctx, endpoint, planData.FlavorId.ValueString(), defaultCloudOperationTimeout)
	if err != nil {
		resp.Diagnostics.AddError("error waiting for load balancer to be ready", err.Error())
		return
	}

	responseData.MergeWith(&planData)

	// Explicitely set the 'name' attribute from the plan as it may
	// take a few seconds before being updated on API side
	responseData.Name = planData.Name

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectLoadbalancerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudProjectRegionLoadbalancerModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/loadbalancing/loadbalancer/" + url.PathEscape(data.Id.ValueString())
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
	}
}

func (r *cloudProjectLoadbalancerResource) waitForLoadBalancerToBeReady(ctx context.Context, path, expectedFlavor string, timeout time.Duration) (*CloudProjectRegionLoadbalancerModel, error) {
	var responseData CloudProjectRegionLoadbalancerModel

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		var resp CloudProjectRegionLoadbalancerModel
		if err := r.config.OVHClient.GetWithContext(ctx, path, &resp); err != nil {
			return retry.NonRetryableError(err)
		}

		if resp.FlavorId.ValueString() != expectedFlavor {
			time.Sleep(5 * time.Second)
			return retry.RetryableError(fmt.Errorf("waiting for load balancer to have the expected flavor (current: %s, expected: %s)", resp.FlavorId.ValueString(), expectedFlavor))
		}

		responseData = resp

		return nil
	})

	return &responseData, err
}

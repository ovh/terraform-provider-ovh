package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/ovh/go-ovh/ovh"
)

var _ resource.ResourceWithConfigure = (*vrackPublicRoutingPriorityResource)(nil)
var _ resource.ResourceWithImportState = (*vrackPublicRoutingPriorityResource)(nil)

func NewVrackPublicRoutingPriorityResource() resource.Resource {
	return &vrackPublicRoutingPriorityResource{}
}

type vrackPublicRoutingPriorityResource struct {
	config *Config
}

// _vrack_public_routing_priority
func (r *vrackPublicRoutingPriorityResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vrack_public_routing_priority"
}

func (d *vrackPublicRoutingPriorityResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *vrackPublicRoutingPriorityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = VrackPublicRoutingPriorityResourceSchema(ctx)
}

func (r *vrackPublicRoutingPriorityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError(
			"Given ID is malformed",
			"import ID must be SERVICE_NAME,PriorityId formatted",
		)
		return
	}

	serviceName := splits[0]
	uuid := splits[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), serviceName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("priority_id"), uuid)...)
}

func (r *vrackPublicRoutingPriorityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, planData, responseData VrackPublicRoutingPriorityModel
	var task VrackTask
	var responseDatas []string

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/publicRoutingPriority"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &task); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	if err := waitForVrackTask(&task, r.config.OVHClient); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error waiting for vrack (%s) to update publicRoutingPriority %v: %s", task.ServiceName, planData.PriorityId, err),
			err.Error(),
		)
		return
	}

	// List vRack Public Routing priorities to fetch the uuid.
	listEndpoint := fmt.Sprintf("/vrack/%s/publicRoutingPriority",
		url.PathEscape(data.ServiceName.ValueString()),
	)
	if err := r.config.OVHClient.Get(listEndpoint, &responseDatas); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", listEndpoint),
			err.Error(),
		)
		return
	}
	for i := range responseDatas {
		priorityID := responseDatas[i]
		endpoint := "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/publicRoutingPriority/" + url.PathEscape(priorityID)
		var respData VrackPublicRoutingPriorityModel
		if err := r.config.OVHClient.Get(endpoint, &respData); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Get %s", listEndpoint),
				err.Error(),
			)
		}
		if respData.Region == data.Region {
			data.PriorityId = respData.PriorityId
			data.Region = respData.Region
			data.Type = respData.Type
			data.Vrack = respData.Vrack
			data.AvailabilityZones = respData.AvailabilityZones
			break
		}
	}

	responseData.MergeWith(&data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *vrackPublicRoutingPriorityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData VrackPublicRoutingPriorityModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/publicRoutingPriority/" + url.PathEscape(data.PriorityId.ValueString())

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

func (r *vrackPublicRoutingPriorityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData, responseData VrackPublicRoutingPriorityModel
	var task VrackTask

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
	endpoint := "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/publicRoutingPriority/" + url.PathEscape(data.PriorityId.ValueString())
	if err := r.config.OVHClient.Put(endpoint, planData.ToUpdate(), &task); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	if err := waitForVrackTask(&task, r.config.OVHClient); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error waiting for vrack (%d) to update publicRoutingPriority %v: %s", task.Id, planData.PriorityId, err),
			err.Error(),
		)
		return
	}

	// Read updated resource
	endpoint = "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/publicRoutingPriority/" + url.PathEscape(data.PriorityId.ValueString())
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

func (r *vrackPublicRoutingPriorityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data, planData VrackPublicRoutingPriorityModel
	var task VrackTask

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/publicRoutingPriority/" + url.PathEscape(data.PriorityId.ValueString())
	if err := r.config.OVHClient.Delete(endpoint, &task); err != nil {
		var apiErr *ovh.APIError
		if errors.As(err, &apiErr) && apiErr.Code == 403 && strings.Contains(apiErr.Message, "cannot be deleted") {
			// Default public routing priority cannot be deleted; remove from state only
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
		return
	}

	if err := waitForVrackTask(&task, r.config.OVHClient); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error waiting for vrack (%s) to update publicRoutingPriority %v: %s", task.ServiceName, planData.PriorityId, err),
			err.Error(),
		)
		return
	}
}

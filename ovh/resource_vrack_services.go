package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/ovh/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ resource.ResourceWithConfigure = (*vrackServicesResource)(nil)

func NewVrackServicesResource() resource.Resource {
	return &vrackServicesResource{}
}

type vrackServicesResource struct {
	config *Config
}

func (r *vrackServicesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vrack_services"
}

func (d *vrackServicesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *vrackServicesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = VrackServicesResourceSchema(ctx)
}

func (r *vrackServicesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData VrackServicesModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create order and wait for service to be delivered
	order := data.ToOrder()
	if err := orderCreate(order, r.config, "vrackServices"); err != nil {
		resp.Diagnostics.AddError("failed to create order", err.Error())
	}

	// Find service name from order
	orderID := order.Order.OrderId.ValueInt64()
	plans := []PlanValue{}
	resp.Diagnostics.Append(data.Plan.ElementsAs(ctx, &plans, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve service name
	serviceName, err := serviceNameFromOrder(r.config.OVHClient, orderID, plans[0].PlanCode.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to retrieve service name", err.Error())
	}
	data.VrackServicesId = ovhtypes.TfStringValue{
		StringValue: basetypes.NewStringValue(serviceName),
	}

	// Update resource
	endpoint := "/v2/vrackServices/resource/" + url.PathEscape(data.VrackServicesId.ValueString())
	if err := r.config.OVHClient.Put(endpoint, data.ToUpdate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Read updated resource
	endpoint = "/v2/vrackServices/resource/" + url.PathEscape(data.VrackServicesId.ValueString())
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vrackServicesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData VrackServicesModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/vrackServices/resource/" + url.PathEscape(data.VrackServicesId.ValueString())
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

func (r *vrackServicesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData, responseData VrackServicesModel

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
	endpoint := "/v2/vrackServices/resource/" + url.PathEscape(data.VrackServicesId.ValueString())
	if err := r.config.OVHClient.Put(endpoint, planData.ToUpdate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Read updated resource
	endpoint = "/v2/vrackServices/resource/" + url.PathEscape(data.VrackServicesId.ValueString())
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

func (r *vrackServicesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VrackServicesModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.VrackServicesId.ValueString()
	serviceID, err := serviceIDFromServiceNameInQuery(r.config.OVHClient, serviceName)
	if err != nil {
		resp.Diagnostics.AddError("failed to get service id to terminate service", err.Error())
		return
	}

	terminate := func() (string, error) {
		log.Printf("[DEBUG] Will terminate service %s", serviceName)
		endpoint := fmt.Sprintf("/services/%d/terminate", serviceID)
		if err := r.config.OVHClient.Post(endpoint, nil, nil); err != nil {
			if errOvh, ok := err.(*ovh.APIError); ok && (errOvh.Code == 404 || errOvh.Code == 460) {
				return "", nil
			}
			return "", fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
		}
		return serviceName, nil
	}

	confirmTerminate := func(token string) error {
		log.Printf("[DEBUG] Will confirm termination of service %s", serviceName)
		endpoint := fmt.Sprintf("/services/%d/terminate/confirm", serviceID)
		if err := r.config.OVHClient.Post(endpoint, &ConfirmTerminationOpts{Token: token}, nil); err != nil {
			return fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
		}
		return nil
	}

	if err := orderDelete(r.config, terminate, confirmTerminate); err != nil {
		resp.Diagnostics.AddError("failed to delete resource", err.Error())
		return
	}
}

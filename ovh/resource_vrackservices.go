package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ resource.ResourceWithConfigure = (*vrackServicesResource)(nil)
var _ resource.ResourceWithImportState = (*vrackServicesResource)(nil)

func NewVrackServicesResource() resource.Resource {
	return &vrackServicesResource{}
}

type vrackServicesResource struct {
	config *Config
}

func (r *vrackServicesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vrackservices"
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

func (d *vrackServicesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func (r *vrackServicesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VrackServicesModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create order and wait for service to be delivered
	order := data.ToOrder()
	if err := orderCreate(order, r.config, "vrackServices", true, defaultOrderTimeout); err != nil {
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
	data.Id = ovhtypes.TfStringValue{
		StringValue: basetypes.NewStringValue(serviceName),
	}

	// Read created resource to get the Checksum for update
	var createdData VrackServicesModel
	endpoint := "/v2/vrackServices/resource/" + url.PathEscape(data.Id.ValueString())
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &createdData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}
	data.Checksum = ovhtypes.TfStringValue{
		StringValue: basetypes.NewStringValue(createdData.Checksum.ValueString()),
	}

	// Update resource
	endpoint = "/v2/vrackServices/resource/" + url.PathEscape(data.Id.ValueString())
	if err := r.config.OVHClient.Put(endpoint, data.ToUpdate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for udpate
	if err := helpers.WaitForAPIv2ResourceStatusReady(ctx, r.config.OVHClient, endpoint); err != nil {
		resp.Diagnostics.AddError("Error waiting for resource to be ready", err.Error())
		return
	}

	// Fetch up-to-date service info
	var responseData VrackServicesModel
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(&responseData, true)

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

	endpoint := "/v2/vrackServices/resource/" + url.PathEscape(data.Id.ValueString())
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData, true)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vrackServicesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var responseData, data, planData VrackServicesModel

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
	endpoint := "/v2/vrackServices/resource/" + url.PathEscape(data.Id.ValueString())

	if !planData.TargetSpec.IsNull() && !planData.TargetSpec.IsUnknown() {
		// Add the checksum with its current value
		updateData := planData.ToUpdate()
		updateData.Checksum = &data.Checksum
		if err := r.config.OVHClient.Put(endpoint, updateData, nil); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Put %s", endpoint),
				err.Error(),
			)
			return
		}
	}

	// Wait for udpate
	if err := helpers.WaitForAPIv2ResourceStatusReady(ctx, r.config.OVHClient, endpoint); err != nil {
		resp.Diagnostics.AddError("Error waiting for resource to be ready", err.Error())
		return
	}

	// Fetch up-to-date service info
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	responseData.MergeWith(&planData, false)

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

	// Check terminate business rules
	if !data.TargetSpec.Subnets.IsNull() && !data.TargetSpec.Subnets.IsUnknown() {
		for _, elem := range data.TargetSpec.Subnets.Elements() {
			subnet := elem.(TargetSpecSubnetsValue)
			if subnet.ServiceEndpoints.IsNull() || subnet.ServiceEndpoints.IsUnknown() {
				continue
			}

			if len(subnet.ServiceEndpoints.Elements()) > 0 {
				resp.Diagnostics.AddError("failed to delete resource",
					"every existing ServiceEndpoints must be deleted before terminate the resource")
				return
			}
		}
	}

	resourceName := data.Id.ValueString()
	serviceID, err := serviceIdFromResourceName(r.config.OVHClient, resourceName)
	if err != nil {
		resp.Diagnostics.AddError("failed to get service id to terminate service", err.Error())
		return
	}

	terminate := func() (string, error) {
		log.Printf("[DEBUG] Will terminate service %s", resourceName)
		endpoint := fmt.Sprintf("/services/%d/terminate", serviceID)
		if err := r.config.OVHClient.Post(endpoint, nil, nil); err != nil {
			if errOvh, ok := err.(*ovh.APIError); ok && (errOvh.Code == 404 || errOvh.Code == 460) {
				return "", nil
			}
			return "", fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
		}
		return resourceName, nil
	}

	confirmTerminate := func(token string) error {
		log.Printf("[DEBUG] Will confirm termination of service %s", resourceName)
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

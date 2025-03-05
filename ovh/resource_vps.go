package ovh

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.ResourceWithConfigure   = (*vpsResource)(nil)
	_ resource.ResourceWithImportState = (*vpsResource)(nil)
)

func NewVpsResource() resource.Resource {
	return &vpsResource{}
}

type vpsResource struct {
	config *Config
}

func (r *vpsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps"
}

func (d *vpsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *vpsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = VpsResourceSchema(ctx)
}

func (r *vpsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Here we force the attribute "plan" to an empty array because it won't be fetched by the Read function.
	// If we don't do this, Terraform always shows a diff on the following plans (null => []), due to the
	// plan modifier RequiresReplace that initializes the attribute to its zero-value (an empty array).
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("plan"), ovhtypes.TfListNestedValue[PlanValue]{
		ListValue: basetypes.NewListValueMust(PlanValue{}.Type(ctx), make([]attr.Value, 0)),
	})...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), req.ID)...)
}

func (r *vpsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VpsModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create order and wait for service to be delivered
	order := data.ToOrder()
	if err := orderCreate(order, r.config, "vps", true); err != nil {
		resp.Diagnostics.AddError("failed to create order", err.Error())
	}

	// Find service name from order
	orderID := order.Order.OrderId.ValueInt64()
	plans := []PlanValue{}
	resp.Diagnostics.Append(data.Plan.ElementsAs(ctx, &plans, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName, err := serviceNameFromOrder(r.config.OVHClient, orderID, plans[0].PlanCode.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to retrieve service name", err.Error())
	}
	data.ServiceName = types.TfStringValue{
		StringValue: basetypes.NewStringValue(serviceName),
	}

	// Update resource
	endpoint := "/vps/" + url.PathEscape(data.ServiceName.ValueString())
	if err := r.config.OVHClient.Put(endpoint, data.ToUpdate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Read updated resource
	responseData, err := r.waitForVPSUpdate(ctx, serviceName, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching updated resource",
			err.Error(),
		)
		return
	}

	data.MergeWith(responseData)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vpsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VpsModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read resource
	responseData, err := r.waitForVPSUpdate(ctx, data.ServiceName.ValueString(), &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching resource",
			err.Error(),
		)
		return
	}

	data.MergeWith(responseData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vpsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData VpsModel

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
	endpoint := "/vps/" + url.PathEscape(data.ServiceName.ValueString())
	if err := r.config.OVHClient.Put(endpoint, planData.ToUpdate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Read updated resource
	responseData, err := r.waitForVPSUpdate(ctx, data.ServiceName.ValueString(), &planData)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error fetching updated resource %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&planData)
	responseData.MergeWith(&data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *vpsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VpsModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()

	terminate := func() (string, error) {
		log.Printf("[DEBUG] Will terminate vps %s", serviceName)
		endpoint := fmt.Sprintf("/vps/%s/terminate", url.PathEscape(serviceName))
		if err := r.config.OVHClient.Post(endpoint, nil, nil); err != nil {
			if errOvh, ok := err.(*ovh.APIError); ok && (errOvh.Code == 404 || errOvh.Code == 460) {
				return "", nil
			}
			return "", fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
		}
		return serviceName, nil
	}

	confirmTerminate := func(token string) error {
		log.Printf("[DEBUG] Will confirm termination of vps %s", serviceName)
		endpoint := fmt.Sprintf("/vps/%s/confirmTermination", url.PathEscape(serviceName))
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

// waitForVPSUpdate fetches the given VPS 20 times in a loop, sleeping 1 min between each fetch. It is done to ensure that
// field `netbootMode` has been updated since it is not done synchronously.
func (r *vpsResource) waitForVPSUpdate(ctx context.Context, serviceName string, planData *VpsModel) (*VpsModel, error) {
	var responseData VpsModel

	err := retry.RetryContext(ctx, 20*time.Minute, func() *retry.RetryError {
		// Read updated resource
		endpoint := "/vps/" + url.PathEscape(serviceName)
		if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
			return retry.NonRetryableError(fmt.Errorf("error calling GET %s", endpoint))
		}

		// Update is succesfull, return
		if responseData.NetbootMode == planData.NetbootMode ||
			planData.NetbootMode.IsNull() ||
			planData.NetbootMode.IsUnknown() {
			return nil
		}

		return retry.RetryableError(errors.New("waiting for netbootMode to have its expected value"))
	})

	if err != nil {
		return nil, err
	}

	return &responseData, nil
}

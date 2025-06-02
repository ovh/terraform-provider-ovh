package ovh

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/ovh/go-ovh/ovh"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ resource.ResourceWithConfigure = (*dedicatedServerResource)(nil)

func NewDedicatedServerResource() resource.Resource {
	return &dedicatedServerResource{}
}

type dedicatedServerResource struct {
	config *Config
}

func (r *dedicatedServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server"
}

func (d *dedicatedServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *dedicatedServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DedicatedServerResourceSchema(ctx)
}

func (d *dedicatedServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.Private.SetKey(ctx, "is_imported", []byte("true"))...)
	resp.Diagnostics.Append(resp.Private.SetKey(ctx, "run_count", []byte("0"))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), req.ID)...)
}

func (r *dedicatedServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData DedicatedServerModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If service_name is not provided, it means dedicated server has to be ordered
	if data.ServiceName.IsNull() || data.ServiceName.IsUnknown() {
		order := data.ToOrder()
		if err := orderCreate(order, r.config, "baremetalServers", false); err != nil {
			resp.Diagnostics.AddError("failed to create order", err.Error())
			return
		}
		orderID := order.Order.OrderId.ValueInt64()
		data.Order = OrderValue{
			state:   attr.ValueStateKnown,
			OrderId: ovhtypes.NewTfInt64Value(orderID),
		}

		// Wait for order to be completed
		_, err := waitOrderCompletionV2(ctx, r.config, orderID)
		if err != nil {
			timeout := &retry.TimeoutError{}
			if errors.As(err, &timeout) {
				// Delivery took too long, just store the order ID and leave.
				// Resource will have to be untainted before next apply (cf: https://discuss.hashicorp.com/t/partial-resource-create-tainted-state/48905).
				resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
				resp.Diagnostics.AddError("still waiting for order to be completed", fmt.Sprintf("Order %d not delivered yet, saving the ID for later", orderID))
			} else {
				// Got a real error, return it and don't save anything in the state
				resp.Diagnostics.AddError("error waiting for order", err.Error())
			}
			return
		}

		// Find service name from order
		r.getServiceName(ctx, &data, orderID, resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Reinstall server if not blocked by configuration
	if err := r.reinstallDedicatedServer(ctx, data.PreventInstallOnCreate.ValueBool(), false, nil, &data); err != nil {
		resp.Diagnostics.AddError("failed to reinstall server", err.Error())
		return
	}

	// Update resource
	r.updateDedicatedServerResource(ctx, nil, &data, &responseData, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		// TODO: not sure we should return here, maybe save the state instead
		return
	}

	responseData.MergeWith(&data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *dedicatedServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData DedicatedServerModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if server has been correctly delivered. If it is not
	// the case, return an error
	if data.ServiceName.IsNull() || data.ServiceName.IsUnknown() {
		orderID := data.Order.OrderId.ValueInt64()

		// Check if server has been delivered, else just exit with error
		_, err := waitOrderCompletionV2(ctx, r.config, orderID)
		if err != nil {
			// If we got a real error, return it. Otherwise, it means
			// service is still being delivered
			timeout := &retry.TimeoutError{}
			if errors.As(err, &timeout) {
				// Delivery took too long, just leave.
				// Resource will have to be untainted before next apply (cf: https://discuss.hashicorp.com/t/partial-resource-create-tainted-state/48905).
				resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
				resp.Diagnostics.AddError("still waiting for order to be completed", fmt.Sprintf("Order %d not delivered yet, saving the ID for later", orderID))
			} else {
				// Got a real error, return it
				resp.Diagnostics.AddError("error waiting for order", err.Error())
			}
			return
		}

		// Find service name from order
		r.getServiceName(ctx, &data, orderID, resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	endpoint := "/dedicated/server/" + url.PathEscape(data.ServiceName.ValueString())
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&data)

	// Check if resource has just been imported. If it is the case, increase
	// the run counter each time we go through this function. The run counter
	// value is used in the Update function to decide if the server should be
	// reinstalled or not.
	isImported, privDiags := req.Private.GetKey(ctx, "is_imported")
	if privDiags.HasError() {
		resp.Diagnostics.Append(privDiags...)
		return
	}

	if isImported != nil && string(isImported) == "true" {
		runCounter, privDiags := req.Private.GetKey(ctx, "run_count")
		if privDiags.HasError() {
			resp.Diagnostics.Append(privDiags...)
			return
		}
		count, err := strconv.Atoi(string(runCounter))
		if err != nil {
			resp.Diagnostics.AddError("failed to read run_count from private state", err.Error())
			return
		}

		resp.Private.SetKey(ctx, "run_count", []byte(strconv.Itoa(count+1)))
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *dedicatedServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateData, planData, responseData DedicatedServerModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if server has been correctly delivered.
	// If it is not the case, return an error.
	if stateData.ServiceName.IsNull() || stateData.ServiceName.IsUnknown() {
		orderID := stateData.Order.OrderId.ValueInt64()

		// Check if server has been delivered, else just exit with error
		_, err := waitOrderCompletionV2(ctx, r.config, orderID)
		if err != nil {
			// If we got a real error, return it. Otherwise, it means
			// service is still being delivered
			timeout := &retry.TimeoutError{}
			if errors.As(err, &timeout) {
				// Delivery took too long, just leave.
				// Resource will have to be untainted before next apply.
				resp.Diagnostics.AddError("still waiting for order to be completed", fmt.Sprintf("Order %d not delivered yet, saving the ID for later", orderID))
			} else {
				// Got a real error, return it
				resp.Diagnostics.AddError("error waiting for order", err.Error())
			}
			return
		}

		// Find service name from order
		r.getServiceName(ctx, &stateData, orderID, resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Check if resource has just been imported, and remove key from private state.
	// We need this information to know if a reinstall should be forced at this point.
	isImported, privDiags := req.Private.GetKey(ctx, "is_imported")
	if privDiags.HasError() {
		resp.Diagnostics.Append(privDiags...)
		return
	}
	resp.Private.SetKey(ctx, "is_imported", nil)

	// Decide if server installation should be blocked. This happens when resource has
	// just been imported and the parameter "prevent_install_on_import" is true.
	var preventReinstall bool
	if isImported != nil && string(isImported) == "true" {
		runCounter, privDiags := req.Private.GetKey(ctx, "run_count")
		if privDiags.HasError() {
			resp.Diagnostics.Append(privDiags...)
			return
		}
		resp.Private.SetKey(ctx, "run_count", nil)

		count, err := strconv.Atoi(string(runCounter))
		if err != nil {
			resp.Diagnostics.AddError("failed to read run_count from private state", err.Error())
			return
		}

		// When "is_imported" is true, check parameter "prevent_install_on_import" to decide
		// if we should allow a server reinstallation.
		// If "run_count" is > 1, it means that Update was not called at import time, so we just
		// ignore this parameter: changing the value of "prevent_install_on_import" after import
		// should not have any effect.
		preventReinstall = count == 1 && planData.PreventInstallOnImport.ValueBool()
	}

	// Reinstall server (if needed and not blocked)
	if err := r.reinstallDedicatedServer(ctx, preventReinstall, false, &stateData, &planData); err != nil {
		resp.Diagnostics.AddError("failed to reinstall server", err.Error())
		return
	}

	// Update resource
	r.updateDedicatedServerResource(ctx, &stateData, &planData, &responseData, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		// TODO: not sure we should return here, maybe save the state instead
		return
	}

	responseData.MergeWith(&planData)
	responseData.MergeWith(&stateData)

	// Explicitely set Customizations/Properties/Storage to what was defined in the plan
	// as we can't determine the right thing to do in MergeWith function
	responseData.Customizations = planData.Customizations
	responseData.Properties = planData.Properties
	responseData.Storage = planData.Storage

	// Same thing for the flags to control reinstallation, set the plan value explicitly
	responseData.PreventInstallOnCreate = planData.PreventInstallOnCreate
	responseData.PreventInstallOnImport = planData.PreventInstallOnImport
	responseData.KeepServiceAfterDestroy = planData.KeepServiceAfterDestroy

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *dedicatedServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DedicatedServerModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.runPreDestroyActions(ctx, &data); err != nil {
		resp.Diagnostics.AddError("failed to run pre-destroy actions", err.Error())
		return
	}

	if data.KeepServiceAfterDestroy.ValueBool() {
		log.Print("Will remove the resource without terminating service")
		return
	}

	serviceName := data.ServiceName.ValueString()
	if serviceName == "" {
		resp.Diagnostics.AddError("cannot terminate resource, missing service_name",
			"service_name field of resource is empty, maybe the server is not yet delivered")
		return
	}

	terminate := func() (string, error) {
		log.Printf("[DEBUG] Will terminate service %s", serviceName)
		endpoint := fmt.Sprintf("/dedicated/server/%s/terminate", url.PathEscape(serviceName))
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
		endpoint := fmt.Sprintf("/dedicated/server/%s/confirmTermination", url.PathEscape(serviceName))
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

func (r *dedicatedServerResource) runPreDestroyActions(ctx context.Context, data *DedicatedServerModel) error {
	if data.RunActionsBeforeDestroy.IsNull() || data.RunActionsBeforeDestroy.IsUnknown() {
		return nil
	}

	elems := data.RunActionsBeforeDestroy.Elements()
	log.Print("Actions to perform before destroy: ", elems)

	for _, elem := range elems {
		action := elem.(ovhtypes.TfStringValue).ValueString()

		switch action {
		case "reinstall_only_os":
			if err := r.reinstallDedicatedServer(ctx, false, true, nil, data); err != nil {
				return fmt.Errorf("failed to reinstall server: %w", err)
			}
		case "reboot_rescue":
			if err := r.rebootRescue(ctx, data); err != nil {
				return fmt.Errorf("failed to reboot server in rescue mode: %w", err)
			}
		default:
			return fmt.Errorf("invalid pre-destroy action given: %q", action)
		}
	}

	return nil
}

func (r *dedicatedServerResource) rebootRescue(ctx context.Context, data *DedicatedServerModel) error {
	// List available boots
	var boots []int
	endpoint := fmt.Sprintf("/dedicated/server/%s/boot?bootType=rescue", data.ServiceName.ValueString())
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &boots); err != nil {
		return fmt.Errorf("failed to fetch boot options: %w", err)
	}

	if len(boots) == 0 {
		return errors.New("no boot found for rescue mode")
	}

	// Update server with boot ID
	endpoint = fmt.Sprintf("/dedicated/server/%s", data.ServiceName.ValueString())
	if err := r.config.OVHClient.PutWithContext(ctx, endpoint, map[string]any{
		"bootId": boots[0],
	}, nil); err != nil {
		return fmt.Errorf("failed to set boot ID %d for server: %w", boots[0], err)
	}

	// Reboot server
	endpoint += "/reboot"
	task := &DedicatedServerTask{}

	if err := r.config.OVHClient.Post(endpoint, nil, task); err != nil {
		return fmt.Errorf("failed to reboot server: %w", err)
	}

	// Wait for reboot to be done
	if err := waitForDedicatedServerTask(data.ServiceName.ValueString(), task, r.config.OVHClient); err != nil {
		return fmt.Errorf("error waiting for reboot task: %w", err)
	}

	return nil
}

func (r *dedicatedServerResource) reinstallDedicatedServer(ctx context.Context, preventReinstall, onlyOS bool, stateData, planData *DedicatedServerModel) error {
	tflog.Debug(ctx, fmt.Sprintf("Prevent server reinstallation: %t", preventReinstall))
	tflog.Debug(ctx, fmt.Sprintf("State data is nil: %t", stateData == nil))

	if preventReinstall {
		tflog.Debug(ctx, "Prevent reinstallation of the server is true")
		return nil
	}

	// Get service name
	serviceName := planData.ServiceName.ValueString()
	if stateData != nil {
		serviceName = stateData.ServiceName.ValueString()
	}

	var shouldReinstall bool

	switch stateData {
	// Check if we should install server right after it has been delivered
	case nil:
		if planData.Os.ValueString() != "" {
			shouldReinstall = true
		}

	// Classical update when resource already exists.
	// Checks state data against plan data to decide if a reinstall should be triggered.
	default:
		if planData.Os.ValueString() != "" &&
			stateData.Os.ValueString() != planData.Os.ValueString() ||
			!stateData.Customizations.Equal(planData.Customizations) ||
			!stateData.Storage.Equal(planData.Storage) ||
			!stateData.Properties.Equal(planData.Properties) {
			shouldReinstall = true
		}
	}

	// Trigger server reinstallation
	if shouldReinstall {
		log.Print("Triggering server reinstallation")

		task := DedicatedServerTask{}
		endpoint := "/dedicated/server/" + url.PathEscape(serviceName) + "/reinstall"
		if err := r.config.OVHClient.Post(endpoint, planData.ToReinstall(onlyOS), &task); err != nil {
			return fmt.Errorf("error calling Post %s", endpoint)
		}

		// Wait for reinstallation completion
		if err := waitForDedicatedServerTask(serviceName, &task, r.config.OVHClient); err != nil {
			return fmt.Errorf("error during server reinstallation: %s", err.Error())
		}
	}

	return nil
}

func (r *dedicatedServerResource) updateDedicatedServerResource(ctx context.Context, stateData, planData, responseData *DedicatedServerModel, diags *diag.Diagnostics) {
	// Get service name
	serviceName := planData.ServiceName.ValueString()
	if stateData != nil {
		serviceName = stateData.ServiceName.ValueString()
	}

	// PUT the resource
	endpoint := "/dedicated/server/" + url.PathEscape(serviceName)
	if err := r.config.OVHClient.Put(endpoint, planData.ToUpdate(stateData), nil); err != nil {
		diags.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Update display name
	if (stateData != nil && stateData.DisplayName.ValueString() != planData.DisplayName.ValueString()) ||
		(stateData == nil && planData.DisplayName.ValueString() != "") {
		newDisplayName := planData.DisplayName.ValueString()
		if err := serviceUpdateDisplayName(ctx, r.config, "dedicated/server", serviceName, newDisplayName); err != nil {
			diags.AddError("failed to update display name", err.Error())
			return
		}
	}

	// Read updated resource
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		diags.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}
}

func (r *dedicatedServerResource) getServiceName(ctx context.Context, data *DedicatedServerModel, orderID int64, diags diag.Diagnostics) {
	// Extract plan from order
	plans := []PlanValue{}
	diags.Append(data.Plan.ElementsAs(ctx, &plans, false)...)
	if diags.HasError() {
		return
	}

	// Retrieve service name
	serviceName, err := serviceNameFromOrder(r.config.OVHClient, orderID, plans[0].PlanCode.ValueString())
	if err != nil {
		diags.AddError("failed to retrieve service name", err.Error())
		return
	}
	data.ServiceName = ovhtypes.TfStringValue{
		StringValue: basetypes.NewStringValue(serviceName),
	}
}

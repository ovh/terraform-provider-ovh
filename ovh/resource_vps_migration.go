package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.ResourceWithConfigure   = (*vpsMigrationResource)(nil)
	_ resource.ResourceWithImportState = (*vpsMigrationResource)(nil)
)

func NewVpsMigrationResource() resource.Resource {
	return &vpsMigrationResource{}
}

type vpsMigrationResource struct {
	config *Config
}

// vpsMigrationAPIResponse mirrors the vps.migration.VPS2020to2025 API model.
type vpsMigrationAPIResponse struct {
	CurrentPlan    string                         `json:"currentPlan"`
	TargetPlan     *string                        `json:"targetPlan,omitempty"`
	AvailablePlans []vpsMigrationAvailablePlanAPI `json:"availablePlans"`
	Status         string                         `json:"status"`
	Date           *string                        `json:"date,omitempty"`
	Position       *int64                         `json:"position,omitempty"`
}

type vpsMigrationAvailablePlanAPI struct {
	PlanCode string `json:"planCode"`
}

type vpsMigrationCreateBody struct {
	Plan string `json:"plan"`
}

type vpsMigrationRescheduleBody struct {
	Date string `json:"date,omitempty"`
}

func (r *vpsMigrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_migration"
}

func (r *vpsMigrationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.config = config
}

func (r *vpsMigrationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = VpsMigrationResourceSchema(ctx)
}

func (r *vpsMigrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// fetchMigration retrieves the current migration record for the given VPS.
func (r *vpsMigrationResource) fetchMigration(serviceName string) (*vpsMigrationAPIResponse, error) {
	endpoint := "/vps/" + url.PathEscape(serviceName) + "/migration2020"
	resp := &vpsMigrationAPIResponse{}
	if err := r.config.OVHClient.Get(endpoint, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// applyAPIResponse copies fields from the API response into the Terraform model.
func applyAPIResponse(ctx context.Context, m *VpsMigrationModel, api *vpsMigrationAPIResponse) {
	m.CurrentPlan = ovhtypes.NewTfStringValue(api.CurrentPlan)
	m.Status = ovhtypes.NewTfStringValue(api.Status)

	if api.TargetPlan != nil {
		m.TargetPlan = ovhtypes.NewTfStringValue(*api.TargetPlan)
	}

	if api.Date != nil {
		m.ScheduledDate = ovhtypes.NewTfStringValue(*api.Date)
	} else {
		m.ScheduledDate = ovhtypes.TfStringValue{StringValue: basetypes.NewStringNull()}
	}

	if api.Position != nil {
		m.Position = ovhtypes.NewTfInt64Value(*api.Position)
	} else {
		m.Position = ovhtypes.NewTfInt64ValueNull()
	}

	elems := make([]attr.Value, 0, len(api.AvailablePlans))
	for _, p := range api.AvailablePlans {
		elems = append(elems, ovhtypes.NewTfStringValue(p.PlanCode))
	}
	listVal, _ := basetypes.NewListValue(ovhtypes.TfStringType{}, elems)
	m.AvailablePlans = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{ListValue: listVal}
}

func (r *vpsMigrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VpsMigrationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()
	targetPlan := data.TargetPlan.ValueString()

	// Confirm a migration is currently possible.
	current, err := r.fetchMigration(serviceName)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get /vps/%s/migration2020", serviceName),
			err.Error(),
		)
		return
	}
	switch current.Status {
	case "available", "toPlan":
		// good
	case "planned", "ongoing", "done":
		resp.Diagnostics.AddError(
			"VPS migration already in progress or completed",
			fmt.Sprintf("Current migration status is %q. Import the existing migration or clean it up first.", current.Status),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"VPS migration not available",
			fmt.Sprintf("Cannot schedule a migration for this VPS (status=%q).", current.Status),
		)
		return
	}

	// Verify the requested plan is offered.
	planAllowed := false
	for _, p := range current.AvailablePlans {
		if p.PlanCode == targetPlan {
			planAllowed = true
			break
		}
	}
	if !planAllowed {
		available := make([]string, 0, len(current.AvailablePlans))
		for _, p := range current.AvailablePlans {
			available = append(available, p.PlanCode)
		}
		resp.Diagnostics.AddError(
			"Invalid target_plan",
			fmt.Sprintf("target_plan %q is not in the available plans list (%v).", targetPlan, available),
		)
		return
	}

	// Enqueue the migration.
	endpoint := "/vps/" + url.PathEscape(serviceName) + "/migration2020"
	body := &vpsMigrationCreateBody{Plan: targetPlan}
	if err := r.config.OVHClient.Post(endpoint, body, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Optionally apply scheduled_date via PUT.
	if !data.ScheduledDate.IsNull() && !data.ScheduledDate.IsUnknown() && data.ScheduledDate.ValueString() != "" {
		put := &vpsMigrationRescheduleBody{Date: data.ScheduledDate.ValueString()}
		if err := r.config.OVHClient.Put(endpoint, put, nil); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Put %s", endpoint),
				err.Error(),
			)
			return
		}
	}

	// Wait until the migration is recorded as planned/ongoing/done (not for completion).
	final, waitErr := r.waitForMigrationScheduled(ctx, serviceName)
	if waitErr != nil {
		resp.Diagnostics.AddError("Error waiting for migration to be scheduled", waitErr.Error())
		return
	}

	applyAPIResponse(ctx, &data, final)
	data.ID = ovhtypes.NewTfStringValue(serviceName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vpsMigrationResource) waitForMigrationScheduled(ctx context.Context, serviceName string) (*vpsMigrationAPIResponse, error) {
	var last *vpsMigrationAPIResponse
	err := retry.RetryContext(ctx, 10*time.Minute, func() *retry.RetryError {
		api, err := r.fetchMigration(serviceName)
		if err != nil {
			return retry.NonRetryableError(err)
		}
		last = api
		switch api.Status {
		case "planned", "ongoing", "done":
			return nil
		}
		return retry.RetryableError(fmt.Errorf("migration status still %q, waiting", api.Status))
	})
	if err != nil {
		return nil, err
	}
	return last, nil
}

func (r *vpsMigrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VpsMigrationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	api, err := r.fetchMigration(data.ServiceName.ValueString())
	if err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get /vps/%s/migration2020", data.ServiceName.ValueString()),
			err.Error(),
		)
		return
	}

	// If no migration is queued anymore (status fell back to available/notAvailable),
	// the resource is effectively gone from Terraform's point of view.
	switch api.Status {
	case "planned", "ongoing", "done", "toPlan":
		// keep
	default:
		resp.State.RemoveResource(ctx)
		return
	}

	applyAPIResponse(ctx, &data, api)
	data.ID = ovhtypes.NewTfStringValue(data.ServiceName.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vpsMigrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData VpsMigrationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := stateData.ServiceName.ValueString()
	endpoint := "/vps/" + url.PathEscape(serviceName) + "/migration2020"

	planChanged := planData.TargetPlan.ValueString() != stateData.TargetPlan.ValueString()
	dateChanged := planData.ScheduledDate.ValueString() != stateData.ScheduledDate.ValueString()

	if planChanged {
		// Target plan changed: cancel the queued migration and enqueue a new one.
		if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
			if errOvh, ok := err.(*ovh.APIError); !ok || errOvh.Code != 404 {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error calling Delete %s to swap target_plan", endpoint),
					err.Error(),
				)
				return
			}
		}

		body := &vpsMigrationCreateBody{Plan: planData.TargetPlan.ValueString()}
		if err := r.config.OVHClient.Post(endpoint, body, nil); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Post %s", endpoint),
				err.Error(),
			)
			return
		}

		if !planData.ScheduledDate.IsNull() && !planData.ScheduledDate.IsUnknown() && planData.ScheduledDate.ValueString() != "" {
			put := &vpsMigrationRescheduleBody{Date: planData.ScheduledDate.ValueString()}
			if err := r.config.OVHClient.Put(endpoint, put, nil); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error calling Put %s", endpoint),
					err.Error(),
				)
				return
			}
		}
	} else if dateChanged {
		// Only the date changed: PUT to reschedule.
		put := &vpsMigrationRescheduleBody{Date: planData.ScheduledDate.ValueString()}
		if err := r.config.OVHClient.Put(endpoint, put, nil); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Put %s", endpoint),
				err.Error(),
			)
			return
		}
	}

	final, err := r.waitForMigrationScheduled(ctx, serviceName)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for migration to be re-scheduled", err.Error())
		return
	}

	applyAPIResponse(ctx, &planData, final)
	planData.ID = ovhtypes.NewTfStringValue(serviceName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (r *vpsMigrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VpsMigrationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()

	// Refresh status: DELETE only succeeds for queued migrations.
	current, err := r.fetchMigration(serviceName)
	if err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get /vps/%s/migration2020", serviceName),
			err.Error(),
		)
		return
	}

	switch current.Status {
	case "ongoing":
		resp.Diagnostics.AddError(
			"Cannot cancel an ongoing VPS migration",
			"The 2020 to 2025 migration is already in progress and cannot be cancelled. Remove the resource from state with `terraform state rm` to stop tracking it.",
		)
		return
	case "done":
		resp.Diagnostics.AddError(
			"Cannot cancel a completed VPS migration",
			"The 2020 to 2025 migration has already completed and cannot be rolled back. Remove the resource from state with `terraform state rm` to stop tracking it.",
		)
		return
	case "available", "notAvailable":
		// nothing queued, nothing to delete
		return
	}

	endpoint := "/vps/" + url.PathEscape(serviceName) + "/migration2020"
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for the queue to clear.
	err = retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
		api, gerr := r.fetchMigration(serviceName)
		if gerr != nil {
			if errOvh, ok := gerr.(*ovh.APIError); ok && errOvh.Code == 404 {
				return nil
			}
			return retry.NonRetryableError(gerr)
		}
		switch api.Status {
		case "available", "notAvailable", "toPlan":
			return nil
		}
		return retry.RetryableError(errors.New("waiting for migration cancellation"))
	})
	if err != nil {
		resp.Diagnostics.AddError("Error verifying migration cancellation", err.Error())
	}
}

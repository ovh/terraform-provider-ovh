package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func VpsMigrationResourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Unique identifier for the resource (same as service_name).",
		},
		"service_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "VPS service name (e.g. vpsXXXX.ovh.net) to migrate from the 2020 to the 2025 generation.",
			MarkdownDescription: "VPS service name (e.g. `vpsXXXX.ovh.net`) to migrate from the 2020 to the 2025 generation.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"target_plan": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Target plan code for the migration. Must be one of available_plans returned by the API.",
			MarkdownDescription: "Target plan code for the migration. Must be one of `available_plans` returned by the API.",
		},
		"scheduled_date": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Optional:            true,
			Computed:            true,
			Description:         "ISO-8601 datetime at which the migration should run. When omitted, the migration is enqueued immediately.",
			MarkdownDescription: "ISO-8601 datetime at which the migration should run. When omitted, the migration is enqueued immediately.",
		},
		"current_plan": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Current (source) plan of the VPS.",
			MarkdownDescription: "Current (source) plan of the VPS.",
		},
		"status": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Current migration status (available, done, notAvailable, ongoing, planned, toPlan).",
			MarkdownDescription: "Current migration status (`available`, `done`, `notAvailable`, `ongoing`, `planned`, `toPlan`).",
		},
		"position": schema.Int64Attribute{
			CustomType:          ovhtypes.TfInt64Type{},
			Computed:            true,
			Description:         "Position of this migration in the queue, when status is planned or toPlan.",
			MarkdownDescription: "Position of this migration in the queue, when status is `planned` or `toPlan`.",
		},
		"available_plans": schema.ListAttribute{
			CustomType:          ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
			Computed:            true,
			Description:         "List of plan codes that this VPS can be migrated to.",
			MarkdownDescription: "List of plan codes that this VPS can be migrated to.",
		},
	}

	return schema.Schema{
		Description: "Schedules the migration of an OVHcloud VPS from the 2020 to the 2025 generation.",
		Attributes:  attrs,
	}
}

type VpsMigrationModel struct {
	ID             ovhtypes.TfStringValue                             `tfsdk:"id" json:"-"`
	ServiceName    ovhtypes.TfStringValue                             `tfsdk:"service_name" json:"-"`
	TargetPlan     ovhtypes.TfStringValue                             `tfsdk:"target_plan" json:"targetPlan"`
	ScheduledDate  ovhtypes.TfStringValue                             `tfsdk:"scheduled_date" json:"date"`
	CurrentPlan    ovhtypes.TfStringValue                             `tfsdk:"current_plan" json:"currentPlan"`
	Status         ovhtypes.TfStringValue                             `tfsdk:"status" json:"status"`
	Position       ovhtypes.TfInt64Value                              `tfsdk:"position" json:"position"`
	AvailablePlans ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] `tfsdk:"available_plans" json:"-"`
}

func (v *VpsMigrationModel) MergeWith(other *VpsMigrationModel) {
	if (v.ID.IsUnknown() || v.ID.IsNull()) && !other.ID.IsUnknown() {
		v.ID = other.ID
	}
	if (v.ServiceName.IsUnknown() || v.ServiceName.IsNull()) && !other.ServiceName.IsUnknown() {
		v.ServiceName = other.ServiceName
	}
	if (v.TargetPlan.IsUnknown() || v.TargetPlan.IsNull()) && !other.TargetPlan.IsUnknown() {
		v.TargetPlan = other.TargetPlan
	}
	if (v.ScheduledDate.IsUnknown() || v.ScheduledDate.IsNull()) && !other.ScheduledDate.IsUnknown() {
		v.ScheduledDate = other.ScheduledDate
	}
	if (v.CurrentPlan.IsUnknown() || v.CurrentPlan.IsNull()) && !other.CurrentPlan.IsUnknown() {
		v.CurrentPlan = other.CurrentPlan
	}
	if (v.Status.IsUnknown() || v.Status.IsNull()) && !other.Status.IsUnknown() {
		v.Status = other.Status
	}
	if (v.Position.IsUnknown() || v.Position.IsNull()) && !other.Position.IsUnknown() {
		v.Position = other.Position
	}
	if (v.AvailablePlans.IsUnknown() || v.AvailablePlans.IsNull()) && !other.AvailablePlans.IsUnknown() {
		v.AvailablePlans = other.AvailablePlans
	}
}

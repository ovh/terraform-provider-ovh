package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func VpsMigrationDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Returns the current 2020 to 2025 migration state for a given VPS, as well as the available target plans.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Unique identifier for the resource (same as service_name).",
			},
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "VPS service name to query.",
				MarkdownDescription: "VPS service name to query.",
			},
			"target_plan": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Target plan currently set on the migration (if any).",
				MarkdownDescription: "Target plan currently set on the migration (if any).",
			},
			"scheduled_date": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "ISO-8601 datetime at which the migration is scheduled to run, if planned.",
				MarkdownDescription: "ISO-8601 datetime at which the migration is scheduled to run, if planned.",
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
				Description:         "Position of this migration in the queue, when applicable.",
				MarkdownDescription: "Position of this migration in the queue, when applicable.",
			},
			"available_plans": schema.ListAttribute{
				CustomType:          ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
				Computed:            true,
				Description:         "List of plan codes that this VPS can be migrated to.",
				MarkdownDescription: "List of plan codes that this VPS can be migrated to.",
			},
		},
	}
}

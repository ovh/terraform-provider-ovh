package ovh

import (
	"context"
	// "fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/ovh/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func OkmsResourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"display_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Optional:            true,
			Computed:            true,
			Description:         "Set the name displayed in Manager for this KMS",
			MarkdownDescription: "Set the name displayed in Manager for this KMS",
		},
		"iam": schema.SingleNestedAttribute{
			Attributes: map[string]schema.Attribute{
				"display_name": schema.StringAttribute{
					CustomType:          ovhtypes.TfStringType{},
					Computed:            true,
					Description:         "Resource display name",
					MarkdownDescription: "Resource display name",
				},
				"id": schema.StringAttribute{
					CustomType:          ovhtypes.TfStringType{},
					Computed:            true,
					Description:         "Unique identifier of the resource",
					MarkdownDescription: "Unique identifier of the resource",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"tags": schema.MapAttribute{
					CustomType:          ovhtypes.NewTfMapNestedType[ovhtypes.TfStringValue](ctx),
					Computed:            true,
					Description:         "Resource tags. Tags that were internally computed are prefixed with ovh:",
					MarkdownDescription: "Resource tags. Tags that were internally computed are prefixed with ovh:",
				},
				"urn": schema.StringAttribute{
					CustomType:          ovhtypes.TfStringType{},
					Computed:            true,
					Description:         "Unique resource name used in policies",
					MarkdownDescription: "Unique resource name used in policies",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			CustomType: IamType{
				ObjectType: types.ObjectType{
					AttrTypes: IamValue{}.AttributeTypes(ctx),
				},
			},
			Computed:            true,
			Description:         "IAM resource metadata",
			MarkdownDescription: "IAM resource metadata",
		},
		"id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "OKMS ID",
			MarkdownDescription: "OKMS ID",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"kmip_endpoint": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "KMS kmip API endpoint",
			MarkdownDescription: "KMS kmip API endpoint",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"ovh_subsidiary": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "OVH subsidiaries",
			MarkdownDescription: "OVH subsidiaries",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"public_ca": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "KMS public CA (Certificate Authority)",
			MarkdownDescription: "KMS public CA (Certificate Authority)",
		},
		"region": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "KMS region",
			MarkdownDescription: "KMS region",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"rest_endpoint": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "KMS rest API endpoint",
			MarkdownDescription: "KMS rest API endpoint",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"swagger_endpoint": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "KMS rest API swagger UI",
			MarkdownDescription: "KMS rest API swagger UI",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}

	return schema.Schema{
		Attributes: attrs,
	}
}

type OkmsModel struct {
	DisplayName     ovhtypes.TfStringValue `tfsdk:"display_name" json:"displayName"`
	Iam             IamValue               `tfsdk:"iam" json:"iam"`
	Id              ovhtypes.TfStringValue `tfsdk:"id" json:"id"`
	KmipEndpoint    ovhtypes.TfStringValue `tfsdk:"kmip_endpoint" json:"kmipEndpoint"`
	PublicCa        ovhtypes.TfStringValue `tfsdk:"public_ca" json:"publicCa"`
	RestEndpoint    ovhtypes.TfStringValue `tfsdk:"rest_endpoint" json:"restEndpoint"`
	SwaggerEndpoint ovhtypes.TfStringValue `tfsdk:"swagger_endpoint" json:"swaggerEndpoint"`
	OvhSubsidiary   ovhtypes.TfStringValue `tfsdk:"ovh_subsidiary" json:"ovhSubsidiary"`
	Region          ovhtypes.TfStringValue `tfsdk:"region" json:"region"`
}

func (v *OkmsModel) ToCreate(ctx context.Context) (*OrderModel, diag.Diagnostics) {
	// Create order configuration
	configuration, diags := NewPlanConfigurationValue(
		PlanConfigurationValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"label": ovhtypes.NewTfStringValue("region"),
			"value": v.Region,
		},
	)

	if diags.HasError() {
		return nil, diags
	}

	plan, diags := NewPlanValue(
		PlanValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"item_id":      ovhtypes.TfInt64Value{Int64Value: basetypes.NewInt64Unknown()},
			"quantity":     ovhtypes.TfInt64Value{Int64Value: basetypes.NewInt64Value(1)},
			"duration":     ovhtypes.NewTfStringValue("P1M"),
			"plan_code":    ovhtypes.NewTfStringValue("okms"),
			"pricing_mode": ovhtypes.NewTfStringValue("default"),
			"configuration": ovhtypes.TfListNestedValue[PlanConfigurationValue]{
				ListValue: basetypes.NewListValueMust(
					configuration.Type(context.Background()),
					[]attr.Value{configuration},
				),
			},
		},
	)

	if diags.HasError() {
		return nil, diags
	}

	order := &OrderModel{
		Order:         OrderValue{},
		OvhSubsidiary: v.OvhSubsidiary,
		Plan: ovhtypes.TfListNestedValue[PlanValue]{
			ListValue: basetypes.NewListValueMust(
				plan.Type(context.Background()),
				[]attr.Value{plan},
			),
		},
	}

	return order, diags
}

func (v *OkmsModel) MergeWith(other *OkmsModel) {

	if (v.DisplayName.IsUnknown() || v.DisplayName.IsNull()) && !other.DisplayName.IsUnknown() {
		v.DisplayName = other.DisplayName
	}

	if (v.Iam.IsUnknown() || v.Iam.IsNull()) && !other.Iam.IsUnknown() {
		v.Iam = other.Iam
	}

	if (v.Id.IsUnknown() || v.Id.IsNull()) && !other.Id.IsUnknown() {
		v.Id = other.Id
	}

	if (v.KmipEndpoint.IsUnknown() || v.KmipEndpoint.IsNull()) && !other.KmipEndpoint.IsUnknown() {
		v.KmipEndpoint = other.KmipEndpoint
	}

	if (v.PublicCa.IsUnknown() || v.PublicCa.IsNull()) && !other.PublicCa.IsUnknown() {
		v.PublicCa = other.PublicCa
	}

	if (v.RestEndpoint.IsUnknown() || v.RestEndpoint.IsNull()) && !other.RestEndpoint.IsUnknown() {
		v.RestEndpoint = other.RestEndpoint
	}

	if (v.SwaggerEndpoint.IsUnknown() || v.SwaggerEndpoint.IsNull()) && !other.SwaggerEndpoint.IsUnknown() {
		v.SwaggerEndpoint = other.SwaggerEndpoint
	}

	if (v.OvhSubsidiary.IsUnknown() || v.OvhSubsidiary.IsNull()) && !other.OvhSubsidiary.IsUnknown() {
		v.OvhSubsidiary = other.OvhSubsidiary
	}
}

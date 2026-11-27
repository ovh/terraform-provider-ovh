package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// Order configurations of a VPS plan. The catalog has others (infrastructure,
// region, vps_install_rtm): they are optional and left to their default value.
const (
	vpsOrderConfigurationDatacenter = "vps_datacenter"
	vpsOrderConfigurationImage      = "vps_os"
)

// vpsOrderCatalog is the part of GET /order/catalog/public/vps used to order a
// VPS.
type vpsOrderCatalog struct {
	Plans []struct {
		PlanCode       string `json:"planCode"`
		Configurations []struct {
			Name   string   `json:"name"`
			Values []string `json:"values"`
		} `json:"configurations"`
	} `json:"plans"`
}

// vpsOrderOffer holds the values a plan accepts in its order configurations:
// the datacenter codes (vps_datacenter) and the image names (vps_os).
type vpsOrderOffer struct {
	Datacenters, Images []string
}

// vpsOrderOfferFor reads the order configurations of a VPS plan in the public
// order catalog of a subsidiary.
func vpsOrderOfferFor(ctx context.Context, client *ovhwrap.Client, subsidiary, planCode string) (*vpsOrderOffer, error) {
	endpoint := "/order/catalog/public/vps?ovhSubsidiary=" + url.QueryEscape(subsidiary)
	var catalog vpsOrderCatalog
	if err := client.GetUnAuthWithContext(ctx, endpoint, &catalog); err != nil {
		return nil, fmt.Errorf("calling GET %s: %w", endpoint, err)
	}

	for _, plan := range catalog.Plans {
		if plan.PlanCode != planCode {
			continue
		}
		offer := &vpsOrderOffer{}
		for _, configuration := range plan.Configurations {
			switch configuration.Name {
			case vpsOrderConfigurationDatacenter:
				offer.Datacenters = configuration.Values
			case vpsOrderConfigurationImage:
				offer.Images = configuration.Values
			}
		}
		return offer, nil
	}
	return nil, fmt.Errorf("plan %q is not in the VPS order catalog of subsidiary %s", planCode, subsidiary)
}

// vpsOrderModel returns the order of one VPS of plan planCode, in datacenter
// datacenterCode, installed with image imageName.
func vpsOrderModel(ctx context.Context, planCode, datacenterCode, imageName, duration, pricingMode, subsidiary string) (*OrderModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	labels := []string{vpsOrderConfigurationDatacenter, vpsOrderConfigurationImage}
	values := []string{datacenterCode, imageName}
	configurations := make([]attr.Value, 0, len(labels))
	for i, label := range labels {
		configuration, d := NewPlanConfigurationValue(
			PlanConfigurationValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"label": ovhtypes.NewTfStringValue(label),
				"value": ovhtypes.NewTfStringValue(values[i]),
			},
		)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		configurations = append(configurations, configuration)
	}

	plan, d := NewPlanValue(
		PlanValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"item_id":      ovhtypes.TfInt64Value{Int64Value: basetypes.NewInt64Unknown()},
			"quantity":     ovhtypes.TfInt64Value{Int64Value: basetypes.NewInt64Value(1)},
			"duration":     ovhtypes.NewTfStringValue(duration),
			"plan_code":    ovhtypes.NewTfStringValue(planCode),
			"pricing_mode": ovhtypes.NewTfStringValue(pricingMode),
			"configuration": ovhtypes.TfListNestedValue[PlanConfigurationValue]{
				ListValue: basetypes.NewListValueMust(PlanConfigurationValue{}.Type(ctx), configurations),
			},
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return &OrderModel{
		OvhSubsidiary: ovhtypes.NewTfStringValue(subsidiary),
		Plan: ovhtypes.TfListNestedValue[PlanValue]{
			ListValue: basetypes.NewListValueMust(PlanValue{}.Type(ctx), []attr.Value{plan}),
		},
	}, diags
}

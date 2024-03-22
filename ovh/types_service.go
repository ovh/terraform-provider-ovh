package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/ovh/terraform-provider-ovh/ovh/types"
)

// ServiceInfos contains the information returned
// by calls to /serviceType/{serviceId}/serviceInfos
type ServiceInfos struct {
	ServiceID int `json:"serviceId"`
}

// Service contains the information returned by
// calls to /services/{serviceId}
type Service struct {
	Billing ServiceBilling `json:"billing"`
}

func (s *Service) ToPlanValue(ctx context.Context) *types.TfListNestedValue[PlanValue] {
	plan := PlanValue{
		PlanCode: types.TfStringValue{
			StringValue: basetypes.NewStringValue(s.Billing.Plan.Code),
		},
		Duration: types.TfStringValue{
			StringValue: basetypes.NewStringValue(s.Billing.Pricing.Duration),
		},
		PricingMode: types.TfStringValue{
			StringValue: basetypes.NewStringValue(s.Billing.Pricing.PricingMode),
		},
		state: attr.ValueStateKnown,
	}

	planValue := types.TfListNestedValue[PlanValue]{ListValue: basetypes.NewListValueMust(PlanValue{}.Type(ctx), []attr.Value{plan})}

	return &planValue
}

type ServiceBilling struct {
	Plan    ServiceBillingPlan    `json:"plan"`
	Pricing ServiceBillingPricing `json:"pricing"`
}

type ServiceBillingPlan struct {
	Code string `json:"code"`
}

type ServiceBillingPricing struct {
	PricingMode string `json:"pricingMode"`
	Duration    string `json:"duration"`
}

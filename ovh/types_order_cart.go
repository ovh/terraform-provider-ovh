package ovh

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type OrderCartCreateOpts struct {
	OvhSubsidiary string  `json:"ovhSubsidiary"`
	Description   *string `json:"description,omitempty"`
	Expire        *string `json:"expire,omitempty"`
}

func (opts *OrderCartCreateOpts) FromResource(d *schema.ResourceData) *OrderCartCreateOpts {
	opts.OvhSubsidiary = strings.ToUpper(d.Get("ovh_subsidiary").(string))
	opts.Description = helpers.GetNilStringPointerFromData(d, "description")
	opts.Expire = helpers.GetNilStringPointerFromData(d, "expire")

	return opts
}

type OrderCartPlanCreateOpts struct {
	CatalogName *string `json:"catalogName,omitempty"`
	Duration    string  `json:"duration"`
	PlanCode    string  `json:"planCode"`
	PricingMode string  `json:"pricingMode"`
	Quantity    int     `json:"quantity"`
}

func (opts *OrderCartPlanCreateOpts) FromResourceWithPath(d *schema.ResourceData, path string) *OrderCartPlanCreateOpts {
	opts.CatalogName = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.catalog_name", path))
	opts.Duration = d.Get(fmt.Sprintf("%s.duration", path)).(string)
	opts.PlanCode = d.Get(fmt.Sprintf("%s.plan_code", path)).(string)
	opts.PricingMode = d.Get(fmt.Sprintf("%s.pricing_mode", path)).(string)
	return opts
}

func (opts *OrderCartPlanCreateOpts) String() string {
	return fmt.Sprintf(
		"planCode: %s, pricingMode: %s, duration: %s, quantity: %d, catalogName: %v",
		opts.PlanCode,
		opts.PricingMode,
		opts.Duration,
		opts.Quantity,
		*opts.CatalogName,
	)
}

type OrderCartPlanOptionsCreateOpts struct {
	CatalogName *string `json:"catalogName,omitempty"`
	Duration    string  `json:"duration"`
	PlanCode    string  `json:"planCode"`
	PricingMode string  `json:"pricingMode"`
	Quantity    int     `json:"quantity"`
	ItemId      int64   `json:"itemId"`
}

func (opts *OrderCartPlanOptionsCreateOpts) FromResourceWithPath(d *schema.ResourceData, path string) *OrderCartPlanOptionsCreateOpts {
	opts.CatalogName = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.catalog_name", path))
	opts.Duration = d.Get(fmt.Sprintf("%s.duration", path)).(string)
	opts.PlanCode = d.Get(fmt.Sprintf("%s.plan_code", path)).(string)
	opts.PricingMode = d.Get(fmt.Sprintf("%s.pricing_mode", path)).(string)
	return opts
}

func (opts *OrderCartPlanOptionsCreateOpts) String() string {
	return fmt.Sprintf(
		"planCode: %s, pricingMode: %s, duration: %s, quantity: %d, itemId: %d, catalogName: %s",
		opts.PlanCode,
		opts.PricingMode,
		opts.Duration,
		opts.Quantity,
		opts.ItemId,
		*opts.CatalogName,
	)
}

type OrderCart struct {
	CartId      string  `json:"cartId"`
	Description *string `json:"description"`
	Expire      *string `json:"expire"`
	Items       []int64 `json:"items"`
	ReadOnly    bool    `json:"readOnly"`
}

func (v OrderCart) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["cart_id"] = v.CartId
	obj["items"] = v.Items
	obj["read_only"] = v.ReadOnly

	if v.Description != nil {
		obj["description"] = *v.Description
	}

	if v.Expire != nil {
		obj["expire"] = *v.Expire
	}

	return obj
}

type OrderCartItemCreateOpts struct {
	CatalogName *string `json:"catalogName,omitempty"`
	Duration    string  `json:"duration"`
	PlanCode    string  `json:"planCode"`
	PricingMode string  `json:"pricingMode"`
	Quantity    int     `json:"quantity"`
}

func (opts *OrderCartItemCreateOpts) FromResource(d *schema.ResourceData) *OrderCartItemCreateOpts {
	opts.CatalogName = helpers.GetNilStringPointerFromData(d, "catalog_name")
	opts.Duration = strings.ToUpper(d.Get("duration").(string))
	opts.PlanCode = strings.ToUpper(d.Get("plan_code").(string))
	opts.PricingMode = strings.ToUpper(d.Get("pricing_mode").(string))
	opts.Quantity = d.Get("quantity").(int)

	return opts
}

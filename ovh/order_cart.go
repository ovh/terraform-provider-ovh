package ovh

import (
	"fmt"
	"log"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

var (
	genericProductPriceSchema = map[string]*schema.Schema{
		"capacities": {
			Type:        schema.TypeList,
			Elem:        schema.TypeString,
			Description: "Capacities of the pricing (type of pricing)",
			Computed:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Description of the pricing",
			Computed:    true,
		},
		"duration": {
			Type:        schema.TypeString,
			Description: "Duration for ordering the product",
			Computed:    true,
		},
		"interval": {
			Type:        schema.TypeInt,
			Description: "Interval of renewal",
			Computed:    true,
		},
		"maximum_quantity": {
			Type:        schema.TypeInt,
			Description: "Maximum quantity that can be ordered",
			Computed:    true,
		},
		"maximum_repeat": {
			Type:        schema.TypeInt,
			Description: "Maximum repeat for renewal",
			Computed:    true,
		},
		"minimum_quantity": {
			Type:        schema.TypeInt,
			Description: "Minimum quantity that can be ordered",
			Computed:    true,
		},
		"minimum_repeat": {
			Type:        schema.TypeInt,
			Description: "Minimum repeat for renewal",
			Computed:    true,
		},
		"price": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Price of the product (Price with its currency and textual representation)",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"currency_code": {
						Type:        schema.TypeString,
						Description: "Currency code",
						Computed:    true,
					},
					"text": {
						Type:        schema.TypeString,
						Description: "Textual representation",
						Computed:    true,
					},
					"value": {
						Type:        schema.TypeFloat,
						Description: "The effective price",
						Computed:    true,
					},
				},
			},
		},
		"price_in_ucents": {
			Type:        schema.TypeInt,
			Description: "Price of the product in micro-centims",
			Computed:    true,
		},
		"pricing_mode": {
			Type:        schema.TypeString,
			Description: "Pricing model identifier",
			Computed:    true,
		},
		"pricing_type": {
			Type:        schema.TypeString,
			Description: "Pricing type",
			Computed:    true,
		},
	}

	genericProductSchema = map[string]*schema.Schema{
		"plan_code": {
			Type:        schema.TypeString,
			Description: "Product offer identifier",
			Computed:    true,
		},
		"product_name": {
			Type:        schema.TypeString,
			Description: "Name of the product",
			Computed:    true,
		},
		"product_type": {
			Type:        schema.TypeString,
			Description: "Product type",
			Computed:    true,
		},

		"prices": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Prices of the product offer",
			Elem: &schema.Resource{
				Schema: genericProductPriceSchema,
			},
		},
	}

	genericOptionsSchema = map[string]*schema.Schema{
		"exclusive": {
			Type:        schema.TypeBool,
			Description: "Define if options of this family are exclusive with each other",
			Computed:    true,
		},
		"family": {
			Type:        schema.TypeString,
			Description: "Option family",
			Computed:    true,
		},
		"mandatory": {
			Type:        schema.TypeBool,
			Description: "Define if an option of this family is mandatory",
			Computed:    true,
		},
	}
)

func orderCartGenericProductSchema() map[string]*schema.Schema {
	schema := map[string]*schema.Schema{
		"cart_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"product": {
			Type:        schema.TypeString,
			Description: "Product",
			Required:    true,
		},

		"result": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of representations of a generic product",
			Elem: &schema.Resource{
				Schema: genericProductSchema,
			},
		},
	}

	return schema
}

func orderCartGenericProductPlanSchema() map[string]*schema.Schema {
	schema := map[string]*schema.Schema{
		"cart_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"product": {
			Type:        schema.TypeString,
			Description: "Product",
			Required:    true,
		},
		"plan_code": {
			Type:     schema.TypeString,
			Required: true,
		},
		"price_capacity": {
			Type:        schema.TypeString,
			Description: "Capacity of the pricing (type of pricing)",
			Required:    true,
		},
		"catalog_name": {
			Type:        schema.TypeString,
			Description: "Catalog name",
			Optional:    true,
		},
		"selected_price": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Selected Price according to capacity",
			Elem: &schema.Resource{
				Schema: genericProductPriceSchema,
			},
		},
	}

	for k, v := range genericProductSchema {
		if k != "plan_code" {
			schema[k] = v
		}
	}

	return schema
}

func orderCartGenericOptionsSchema() map[string]*schema.Schema {
	resultSchemaAttrs := map[string]*schema.Schema{}
	for k, v := range genericProductSchema {
		resultSchemaAttrs[k] = v
	}

	for k, v := range genericOptionsSchema {
		resultSchemaAttrs[k] = v
	}

	schema := map[string]*schema.Schema{
		"cart_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"plan_code": {
			Type:        schema.TypeString,
			Description: "Product offer identifier",
			Required:    true,
		},
		"product": {
			Type:        schema.TypeString,
			Description: "Product",
			Required:    true,
		},
		"catalog_name": {
			Type:        schema.TypeString,
			Description: "Catalog name",
			Optional:    true,
		},

		"result": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of representations of a generic product",
			Elem: &schema.Resource{
				Schema: resultSchemaAttrs,
			},
		},
	}

	return schema
}

func orderCartGenericOptionsPlanSchema() map[string]*schema.Schema {
	schema := map[string]*schema.Schema{
		"cart_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"plan_code": {
			Type:     schema.TypeString,
			Required: true,
		},
		"options_plan_code": {
			Type:     schema.TypeString,
			Required: true,
		},
		"price_capacity": {
			Type:        schema.TypeString,
			Description: "Capacity of the pricing (type of pricing)",
			Required:    true,
		},
		"product": {
			Type:        schema.TypeString,
			Description: "Product",
			Required:    true,
		},
		"catalog_name": {
			Type:        schema.TypeString,
			Description: "Catalog name",
			Optional:    true,
		},

		"selected_price": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Selected Price according to capacity",
			Elem: &schema.Resource{
				Schema: genericProductPriceSchema,
			},
		},
	}

	for k, v := range genericProductSchema {
		if k != "plan_code" {
			schema[k] = v
		}
	}

	for k, v := range genericOptionsSchema {
		schema[k] = v
	}

	return schema
}

func orderCartGenericProductRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cartId := d.Get("cart_id").(string)
	product := d.Get("product").(string)

	log.Printf("[DEBUG] Will read order cart %s for cart: %s", product, cartId)

	res := []OrderCartGenericProduct{}

	endpoint := fmt.Sprintf(
		"/order/cart/%s/%s",
		url.PathEscape(cartId),
		product,
	)

	err := config.OVHClient.Get(endpoint, &res)
	if err != nil {
		return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
	}

	plans := make([]map[string]interface{}, len(res))
	codes := make([]string, len(res))

	for i, plan := range res {
		plans[i] = plan.ToMap()
		codes[i] = plan.PlanCode
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(codes)

	d.SetId(hashcode.Strings(codes))
	d.Set("result", plans)

	return nil
}

func orderCartGenericOptionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cartId := d.Get("cart_id").(string)
	planCode := d.Get("plan_code").(string)
	product := d.Get("product").(string)

	log.Printf("[DEBUG] Will read order cart options %s for cart: %s", product, cartId)

	res := []OrderCartGenericOptions{}

	endpoint := fmt.Sprintf(
		"/order/cart/%s/%s/options?planCode=%s",
		url.PathEscape(cartId),
		product,
		url.PathEscape(planCode),
	)

	err := config.OVHClient.Get(endpoint, &res)
	if err != nil {
		return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
	}

	plans := make([]map[string]interface{}, len(res))
	codes := make([]string, len(res))

	for i, plan := range res {
		plans[i] = plan.ToMap()
		codes[i] = plan.PlanCode
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(codes)

	d.SetId(hashcode.Strings(codes))
	d.Set("result", plans)

	return nil
}

func orderCartGenericProductPlanRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cartId := d.Get("cart_id").(string)
	planCode := d.Get("plan_code").(string)
	priceCapacity := d.Get("price_capacity").(string)
	product := d.Get("product").(string)

	log.Printf("[DEBUG] Will read order cart %s for cart: %s", product, cartId)

	res := []OrderCartGenericProduct{}

	endpoint := fmt.Sprintf(
		"/order/cart/%s/%s",
		url.PathEscape(cartId),
		product,
	)

	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
	}

	match := false
	matchPrice := false
	for _, plan := range res {
		if plan.PlanCode == planCode {
			match = true
			for k, v := range plan.ToMap() {
				if k != "id" {
					d.Set(k, v)
				}
			}

			// find Price
			for _, price := range plan.Prices {
				for _, cap := range price.Capacities {
					if cap == priceCapacity {
						matchPrice = true
						d.Set("selected_price", []interface{}{price.ToMap()})
						break
					}
				}
			}
			if matchPrice {
				d.SetId(planCode)
			}
		}
	}

	if !match || !matchPrice {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	return nil
}

func orderCartGenericOptionsPlanRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cartId := d.Get("cart_id").(string)
	planCode := d.Get("plan_code").(string)
	optionsPlanCode := d.Get("options_plan_code").(string)
	priceCapacity := d.Get("price_capacity").(string)
	product := d.Get("product").(string)

	log.Printf("[DEBUG] Will read order cart %s for cart: %s", product, cartId)

	res := []OrderCartGenericOptions{}

	endpoint := fmt.Sprintf(
		"/order/cart/%s/%s/options?planCode=%s",
		url.PathEscape(cartId),
		product,
		url.PathEscape(planCode),
	)

	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
	}

	match := false
	matchPrice := false
	for _, plan := range res {
		if plan.PlanCode == optionsPlanCode {
			match = true
			for k, v := range plan.ToMap() {
				if k != "id" {
					d.Set(k, v)
				}
			}

			// find Price
			for _, price := range plan.Prices {
				for _, cap := range price.Capacities {
					if cap == priceCapacity {
						matchPrice = true
						d.Set("selected_price", []interface{}{price.ToMap()})
						break
					}
				}
			}
			if matchPrice {
				d.SetId(optionsPlanCode)
			}
		}
	}

	if !match || !matchPrice {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	return nil
}

func orderCartCreate(meta interface{}, params *OrderCartCreateOpts, assign bool) (*OrderCart, error) {
	config := meta.(*Config)
	r := &OrderCart{}

	log.Printf("[DEBUG] Will create order cart: %v", params)
	endpoint := fmt.Sprintf(
		"/order/cart",
	)

	err := config.OVHClient.Post(endpoint, params, r)
	if err != nil {
		return nil, fmt.Errorf("calling Post %s with params %v:\n\t %q", endpoint, params, err)
	}

	if assign {
		log.Printf("[DEBUG] Will assign order cart: %v", params)
		assign_endpoint := fmt.Sprintf(
			"/order/cart/%s/assign",
			url.PathEscape(r.CartId),
		)

		err = config.OVHClient.Post(assign_endpoint, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("calling Post %s:\n\t %q", assign_endpoint, err)
		}
	}

	return r, nil
}

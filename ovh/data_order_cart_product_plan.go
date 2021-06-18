package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOrderCartProductPlan() *schema.Resource {
	return &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error {
			return orderCartGenericProductPlanRead(d, meta)
		},
		Schema: orderCartGenericProductPlanSchema(),
	}
}

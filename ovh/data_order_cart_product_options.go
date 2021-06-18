package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOrderCartProductOptions() *schema.Resource {
	return &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error {
			return orderCartGenericOptionsRead(d, meta)
		},
		Schema: orderCartGenericOptionsSchema(),
	}
}

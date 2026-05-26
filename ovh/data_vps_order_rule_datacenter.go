package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPSOrderRuleDatacenter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSOrderRuleDatacenterRead,
		Schema: map[string]*schema.Schema{
			"ovh_subsidiary": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "OVH subsidiary (e.g. FR, US, CA, DE, GB).",
			},
			"plan_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Plan code of the VPS to order.",
			},
			"os": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Operating system filter.",
			},
			"datacenters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code":                 {Type: schema.TypeString, Computed: true},
						"datacenter":           {Type: schema.TypeString, Computed: true},
						"days_before_delivery": {Type: schema.TypeInt, Computed: true},
						"status":               {Type: schema.TypeString, Computed: true},
						"linux_status":         {Type: schema.TypeString, Computed: true},
						"windows_status":       {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceVPSOrderRuleDatacenterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	subsidiary := d.Get("ovh_subsidiary").(string)
	planCode := d.Get("plan_code").(string)

	params := url.Values{}
	params.Set("ovhSubsidiary", subsidiary)
	params.Set("planCode", planCode)
	if os, ok := d.GetOk("os"); ok {
		params.Set("os", os.(string))
	}
	endpoint := fmt.Sprintf("/vps/order/rule/datacenter?%s", params.Encode())

	res := VPSOrderRuleDatacenters{}
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return fmt.Errorf("Error calling GET %s: %w", endpoint, err)
	}

	out := make([]map[string]interface{}, 0, len(res.Datacenters))
	for _, dc := range res.Datacenters {
		out = append(out, map[string]interface{}{
			"code":                 dc.Code,
			"datacenter":           dc.Datacenter,
			"days_before_delivery": dc.DaysBeforeDelivery,
			"status":               dc.Status,
			"linux_status":         dc.LinuxStatus,
			"windows_status":       dc.WindowsStatus,
		})
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", subsidiary, planCode, d.Get("os").(string)))
	d.Set("datacenters", out)
	return nil
}

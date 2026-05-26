package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPSOrderRuleOSChoices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSOrderRuleOSChoicesRead,
		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Datacenter code (e.g. GRA, BHS, SBG).",
			},
			"os": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Operating system family.",
			},
			"choices": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name":   {Type: schema.TypeString, Computed: true},
						"status": {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceVPSOrderRuleOSChoicesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	dc := d.Get("datacenter").(string)
	os := d.Get("os").(string)

	params := url.Values{}
	params.Set("datacenter", dc)
	params.Set("os", os)
	endpoint := fmt.Sprintf("/vps/order/rule/osChoices?%s", params.Encode())

	res := VPSOrderRuleOSChoices{}
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return fmt.Errorf("Error calling GET %s: %w", endpoint, err)
	}

	out := make([]map[string]interface{}, 0, len(res.Choices))
	for _, c := range res.Choices {
		out = append(out, map[string]interface{}{
			"name":   c.Name,
			"status": c.Status,
		})
	}

	d.SetId(fmt.Sprintf("%s/%s", dc, os))
	d.Set("choices", out)
	return nil
}

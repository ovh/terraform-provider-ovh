package ovh

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func dataSourceIpLoadbalancingVrackNetworks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIpLoadbalancingVrackNetworksRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The internal name of your iploadbalancer.",
				Required:    true,
			},

			"subnet": {
				Type:        schema.TypeString,
				Description: "Filters on subnet",
				Optional:    true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateIpBlock(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"vlan_id": {
				Type:        schema.TypeInt,
				Description: "Filters on vlan id",
				Optional:    true,
			},

			"result": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourceIpLoadbalancingVrackNetworksRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	result := make([]int64, 0)

	serviceName := d.Get("service_name").(string)
	vlanId := ""
	subnet := ""
	filters := ""

	endpoint := fmt.Sprintf(
		"/ipLoadbalancing/%s/vrack/network",
		url.PathEscape(serviceName),
	)

	if val, ok := d.GetOkExists("vlan_id"); ok {
		vlanId = strconv.Itoa(val.(int))
		filters = fmt.Sprintf("%s&vlan=%s", filters, url.PathEscape(vlanId))
	}

	if val, ok := d.GetOkExists("subnet"); ok {
		subnet = val.(string)
		filters = fmt.Sprintf("%s&subnet=%s", filters, url.PathEscape(subnet))
	}

	// OVH IPLB API doens't parse the query string according to
	// the RFC and throws a 400 error with empty query string
	if filters != "" {
		endpoint = fmt.Sprintf("%s?%s", endpoint, filters)
	}

	if err := config.OVHClient.Get(endpoint, &result); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	d.SetId(fmt.Sprintf("%s_%s_%s", serviceName, subnet, vlanId))
	d.Set("result", result)
	return nil
}

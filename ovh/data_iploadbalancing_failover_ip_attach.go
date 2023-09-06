package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func dataSourceIpLoadbalancingFailoverIpAttach() *schema.Resource {
	return &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error {
			return dataSourceIpLoadbalancingFailoverIpAttachRead(d, meta)
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateIpBlock(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
		},
	}
}

func dataSourceIpLoadbalancingFailoverIpAttachRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/failover", url.PathEscape(serviceName))

	ipBlocks := []string{}
	if err := config.OVHClient.Get(endpoint, &ipBlocks); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	match := false
	for _, ip := range ipBlocks {
		if ip == d.Get("ip").(string) {
			match = true
			d.SetId(serviceName)
		}
	}

	if !match {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return nil
}

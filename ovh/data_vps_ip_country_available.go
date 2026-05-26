package ovh

import (
	"fmt"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

func dataSourceVPSIpCountryAvailable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSIpCountryAvailableRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"result": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceVPSIpCountryAvailableRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	countries := []string{}
	endpoint := fmt.Sprintf(
		"/vps/%s/ipCountryAvailable",
		url.PathEscape(serviceName),
	)
	if err := config.OVHClient.Get(endpoint, &countries); err != nil {
		return fmt.Errorf("Error calling GET %s: %q", endpoint, err)
	}

	sort.Strings(countries)
	d.SetId(hashcode.Strings(append([]string{serviceName}, countries...)))
	d.Set("result", countries)
	return nil
}

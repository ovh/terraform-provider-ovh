package ovh

import (
	"fmt"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

func dataSourceVPSSecondaryDNSDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSSecondaryDNSDomainsRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"result": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVPSSecondaryDNSDomainsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/vps/%s/secondaryDnsDomains", url.PathEscape(serviceName))
	domains := []string{}
	if err := config.OVHClient.Get(endpoint, &domains); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t%q", endpoint, err)
	}

	sort.Strings(domains)
	d.SetId(hashcode.Strings(append([]string{serviceName}, domains...)))
	d.Set("result", domains)
	return nil
}

package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPSSecondaryDNSNameServerAvailable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSSecondaryDNSNameServerAvailableRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVPSSecondaryDNSNameServerAvailableRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/vps/%s/secondaryDnsNameServerAvailable", url.PathEscape(serviceName))
	resp := &VPSSecondaryDNSNameServer{}
	if err := config.OVHClient.Get(endpoint, resp); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t%q", endpoint, err)
	}

	d.SetId(serviceName + "|" + resp.Hostname)
	d.Set("hostname", resp.Hostname)
	d.Set("ip", resp.Ip)
	d.Set("ipv6", resp.Ipv6)
	return nil
}

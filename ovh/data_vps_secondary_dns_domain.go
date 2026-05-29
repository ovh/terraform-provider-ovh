package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPSSecondaryDNSDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSSecondaryDNSDomainRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dns": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_master": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVPSSecondaryDNSDomainRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	domain := d.Get("domain").(string)

	endpoint := fmt.Sprintf("/vps/%s/secondaryDnsDomains/%s",
		url.PathEscape(serviceName), url.PathEscape(domain))
	resp := &VPSSecondaryDNSDomain{}
	if err := config.OVHClient.Get(endpoint, resp); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t%q", endpoint, err)
	}

	d.SetId(serviceName + "|" + resp.Domain)
	d.Set("dns", resp.Dns)
	d.Set("ip_master", resp.IpMaster)
	d.Set("creation_date", resp.CreationDate)
	return nil
}

package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPSIp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSIpRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Computed
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"geolocation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reverse": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVPSIpRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	ipAddress := d.Get("ip_address").(string)

	ip := &VPSIp{}
	endpoint := fmt.Sprintf(
		"/vps/%s/ips/%s",
		url.PathEscape(serviceName),
		url.PathEscape(ipAddress),
	)
	if err := config.OVHClient.Get(endpoint, ip); err != nil {
		return fmt.Errorf("Error calling GET %s: %q", endpoint, err)
	}

	d.SetId(serviceName + "|" + ipAddress)
	d.Set("version", ip.Version)
	d.Set("type", ip.Type)
	if ip.Gateway != nil {
		d.Set("gateway", *ip.Gateway)
	}
	if ip.MacAddress != nil {
		d.Set("mac_address", *ip.MacAddress)
	}
	d.Set("geolocation", ip.Geolocation)
	if ip.Reverse != nil {
		d.Set("reverse", *ip.Reverse)
	}
	return nil
}

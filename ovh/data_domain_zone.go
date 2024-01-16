package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDomainZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDomainZoneRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"has_dns_anycast": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"dnssec_supported": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name_servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"last_update": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDomainZoneRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	zoneName := d.Get("name").(string)

	dz := &DomainZone{}
	err := config.OVHClient.Get(fmt.Sprintf("/domain/zone/%s", url.PathEscape(zoneName)), &dz)

	if err != nil {
		return fmt.Errorf("Error calling /domain/zone/%s:\n\t %q", zoneName, err)
	}

	d.SetId(zoneName)
	d.Set("has_dns_anycast", dz.HasDnsAnycast)
	d.Set("dnssec_supported", dz.DnssecSupported)
	d.Set("last_update", dz.LastUpdate)
	d.Set("name_servers", dz.NameServers)
	d.Set("urn", dz.URN)

	return nil
}

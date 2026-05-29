package ovh

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

func dataSourceVPSDatacenters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSDatacentersRead,
		Schema: map[string]*schema.Schema{
			"country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter datacenters by country code (e.g. FR, US, CA, DE).",
			},
			"datacenters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVPSDatacentersRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	endpoint := "/vps/datacenter"
	if c, ok := d.GetOk("country"); ok {
		endpoint = fmt.Sprintf("%s?country=%s", endpoint, url.QueryEscape(c.(string)))
	}

	dcs := []string{}
	if err := config.OVHClient.Get(endpoint, &dcs); err != nil {
		if apiErr, ok := err.(*ovh.APIError); ok && apiErr.Code == 404 {
			msg := apiErr.Message
			switch {
			case strings.Contains(msg, "Got an invalid (or empty) URL"):
				return fmt.Errorf(
					"the OVHcloud API endpoint %s is not available on this VPS lineup. "+
						"This data source may only work on legacy VPS plans, or the endpoint "+
						"may have been deprecated. See the data source's documentation for "+
						"supported VPS generations.",
					endpoint)
			case strings.Contains(msg, "does not exist"):
				return fmt.Errorf(
					"the requested resource at %s does not exist (the VPS may not have "+
						"the required option subscribed, or the resource ID is wrong)",
					endpoint)
			}
		}
		return fmt.Errorf("calling GET %s: %w", endpoint, err)
	}

	sort.Strings(dcs)
	d.SetId(hashcode.Strings(dcs))
	d.Set("datacenters", dcs)
	return nil
}

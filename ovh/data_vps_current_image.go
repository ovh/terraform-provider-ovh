package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

func dataSourceVPSCurrentImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSCurrentImageRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVPSCurrentImageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	img := &VPSImage{}
	endpoint := fmt.Sprintf(
		"/vps/%s/images/current",
		url.PathEscape(serviceName),
	)
	if err := config.OVHClient.Get(endpoint, img); err != nil {
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

	d.SetId(fmt.Sprintf("%s/%s", serviceName, img.ID))
	d.Set("id", img.ID)
	d.Set("name", img.Name)
	return nil
}

package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		return fmt.Errorf("Error calling GET %s: %w", endpoint, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", serviceName, img.ID))
	d.Set("id", img.ID)
	d.Set("name", img.Name)
	return nil
}

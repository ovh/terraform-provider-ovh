package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// VPSImage represents a VPS image as returned by /vps/{sn}/images/available/{id}
// and /vps/{sn}/images/current.
type VPSImage struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func dataSourceVPSAvailableImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSAvailableImageRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVPSAvailableImageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	imageID := d.Get("image_id").(string)

	img := &VPSImage{}
	endpoint := fmt.Sprintf(
		"/vps/%s/images/available/%s",
		url.PathEscape(serviceName),
		url.PathEscape(imageID),
	)
	if err := config.OVHClient.Get(endpoint, img); err != nil {
		return fmt.Errorf("Error calling GET %s: %w", endpoint, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", serviceName, img.ID))
	d.Set("image_id", img.ID)
	d.Set("name", img.Name)
	return nil
}

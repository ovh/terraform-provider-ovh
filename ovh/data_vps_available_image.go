package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
		},
	}
}

func dataSourceVPSAvailableImageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	imageId := d.Get("image_id").(string)
	image := &VPSImage{}
	err := config.OVHClient.Get(
		fmt.Sprintf(
			"/vps/%s/images/available/%s",
			url.PathEscape(serviceName),
			url.PathEscape(imageId),
		),
		&image,
	)

	if err != nil {
		return nil
	}

	d.Set("image", image)

	return nil
}

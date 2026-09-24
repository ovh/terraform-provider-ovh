package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPSAvailableImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSAvailableImagesRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image_ids": {
				Type:     schema.TypeList,
				Computed: true,
			},
		},
	}
}

func dataSourceVPSAvailableImagesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	imageIds := make([]string, 0)
	err := config.OVHClient.Get(
		fmt.Sprintf(
			"/vps/%s/images/available",
			url.PathEscape(serviceName),
		),
		&imageIds,
	)

	if err != nil {
		return nil
	}

	d.Set("image_ids", imageIds)

	return nil
}

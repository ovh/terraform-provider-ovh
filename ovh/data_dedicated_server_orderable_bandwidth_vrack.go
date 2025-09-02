package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDedicatedServerOrderableBandwidthVrack() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDedicatedServerOrderableBandwidthVrackRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"orderable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is bandwidth orderable for this server",
			},
			"vrack": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Vrack orderable bandwidth in mbps",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func dataSourceDedicatedServerOrderableBandwidthVrackRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	var ds DedicatedServerOrderableBandwidthVrack
	err := config.OVHClient.Get(
		fmt.Sprintf(
			"/dedicated/server/%s/orderable/bandwidthvRack",
			url.PathEscape(serviceName),
		),
		&ds,
	)

	if err != nil {
		return fmt.Errorf(
			"Error calling /dedicated/server/%s/orderable/bandwidthvRack:\n\t %q",
			serviceName,
			err,
		)
	}

	d.SetId(serviceName)
	d.Set("orderable", ds.Orderable)
	d.Set("vrack", ds.Vrack)

	return nil
}

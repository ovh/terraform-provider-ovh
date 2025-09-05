package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDedicatedServerOrderableBandwidth() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDedicatedServerOrderableBandwidthRead,
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
			"platinium": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Platinium orderable bandwidth in mbps",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"premium": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Premium orderable bandwidth in mbps",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"ultimate": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Ultimate orderable bandwidth in mbps",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func dataSourceDedicatedServerOrderableBandwidthRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	var ds DedicatedServerOrderableBandwidth
	err := config.OVHClient.Get(
		fmt.Sprintf(
			"/dedicated/server/%s/orderable/bandwidth",
			url.PathEscape(serviceName),
		),
		&ds,
	)

	if err != nil {
		return fmt.Errorf(
			"Error calling /dedicated/server/%s/orderable/bandwidth:\n\t %q",
			serviceName,
			err,
		)
	}

	d.SetId(serviceName)
	d.Set("orderable", ds.Orderable)
	d.Set("platinium", ds.Platinium)
	d.Set("ultimate", ds.Ultimate)
	d.Set("premium", ds.Premium)

	return nil
}

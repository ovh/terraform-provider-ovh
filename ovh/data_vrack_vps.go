package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVrackVps() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVrackVpsRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Service name of the vrack resource.",
			},
			"vps_service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Service name of the VPS attached to the vRack.",
			},
		},
	}
}

func dataSourceVrackVpsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	vpsServiceName := d.Get("vps_service_name").(string)

	vv := &VrackVps{}
	endpoint := fmt.Sprintf("/vrack/%s/vps/%s",
		url.PathEscape(serviceName),
		url.PathEscape(vpsServiceName),
	)

	if err := config.OVHClient.Get(endpoint, vv); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	d.SetId(fmt.Sprintf("vrack_%s-vps_%s", vv.Vrack, vv.Vps))
	d.Set("service_name", vv.Vrack)
	d.Set("vps_service_name", vv.Vps)
	return nil
}

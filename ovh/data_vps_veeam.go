package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPSVeeam() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSVeeamRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your VPS",
			},
			"backup": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the Veeam backup option is enabled on this VPS",
			},
		},
	}
}

func dataSourceVPSVeeamRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/vps/%s/veeam", url.PathEscape(serviceName))
	veeam := &VpsVeeam{}
	if err := config.OVHClient.Get(endpoint, veeam); err != nil {
		return fmt.Errorf("error calling GET %s: %s", endpoint, err)
	}

	d.SetId(serviceName)
	d.Set("backup", veeam.Backup)
	return nil
}

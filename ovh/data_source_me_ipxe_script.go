package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMeIpxeScript() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMeIpxeScriptRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of your script",
			},
			"script": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Content of your IPXE script",
			},
		},
	}
}

// Common function with the datasource
func dataSourceMeIpxeScriptRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ipxeScript := &MeIpxeScriptResponse{}

	name := d.Get("name").(string)
	err := config.OVHClient.Get(
		fmt.Sprintf("/me/ipxeScript/%s", url.PathEscape(name)),
		ipxeScript,
	)
	if err != nil {
		return fmt.Errorf("Unable to find IpxeScript named %s:\n\t %q", name, err)
	}

	d.SetId(ipxeScript.Name)
	d.Set("name", ipxeScript.Name)
	d.Set("script", ipxeScript.Script)
	return nil
}

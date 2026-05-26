package ovh

import (
	"fmt"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

func dataSourceVrackVpss() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVrackVpssRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Service name of the vrack resource.",
			},
			"result": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVrackVpssRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)

	result := make([]string, 0)
	endpoint := fmt.Sprintf("/vrack/%s/vps", url.PathEscape(serviceName))

	if err := config.OVHClient.Get(endpoint, &result); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	sort.Strings(result)
	d.SetId(hashcode.Strings(append([]string{serviceName}, result...)))
	d.Set("result", result)
	return nil
}

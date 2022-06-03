package ovh

import (
	"fmt"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceVrackCloudProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVrackCloudProjectRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Here come all the computed items
			"result": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVrackCloudProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	result := make([]string, 0)
	err := config.OVHClient.Get(
		fmt.Sprintf(
			"/vrack/%s/cloudProject",
			url.PathEscape(serviceName),
		),
		&result,
	)

	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error calling /vrack/%s/cloudProject:\n\t %q", url.PathEscape(serviceName), err)
	}

	sort.Strings(result)
	d.SetId(hashcode.Strings(result))
	d.Set("result", result)
	return nil
}

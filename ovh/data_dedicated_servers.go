package ovh

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceDedicatedServers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDedicatedServersRead,
		Schema: map[string]*schema.Schema{
			// Computed
			"result": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceDedicatedServersRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ids := []string{}
	err := config.OVHClient.Get("/dedicated/server", &ids)

	if err != nil {
		return fmt.Errorf("Error calling /dedicated/server:\n\t %q", err)
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(ids)

	d.SetId(hashcode.Strings(ids))
	d.Set("result", ids)
	return nil
}

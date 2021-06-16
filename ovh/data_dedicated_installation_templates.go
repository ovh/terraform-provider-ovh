package ovh

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceDedicatedInstallationTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDedicatedInstallationTemplatesRead,
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

func dataSourceDedicatedInstallationTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ids := []string{}
	err := config.OVHClient.Get("/dedicated/installationTemplate", &ids)

	if err != nil {
		return fmt.Errorf("Error calling /dedicated/installationTemplate:\n\t %q", err)
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(ids)

	d.SetId(hashcode.Strings(ids))
	d.Set("result", ids)
	return nil
}

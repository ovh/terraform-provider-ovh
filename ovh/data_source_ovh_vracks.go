package ovh

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceVracks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVracksRead,
		Schema: map[string]*schema.Schema{
			"result": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVracksRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	result := make([]string, 0)
	err := config.OVHClient.Get("/vrack", &result)

	if err != nil {
		return fmt.Errorf("Error calling /vrack:\n\t %q", err)
	}

	sort.Strings(result)
	d.SetId(hashcode.Strings(result))
	d.Set("result", result)
	return nil
}

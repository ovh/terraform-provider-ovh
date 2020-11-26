package ovh

import (
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceMeSshKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMeSshKeysRead,
		Schema: map[string]*schema.Schema{
			"names": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceMeSshKeysRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	names := make([]string, 0)
	err := config.OVHClient.Get("/me/sshKey", &names)

	if err != nil {
		return fmt.Errorf("Error calling /me/sshKey:\n\t %q", err)
	}

	sort.Strings(names)
	d.SetId(hashcode.Strings(names))
	d.Set("names", names)

	log.Printf("[DEBUG] Read SSH Keys names %s", names)
	return nil
}

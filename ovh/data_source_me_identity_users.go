package ovh

import (
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceMeIdentityUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMeIdentityUsersRead,
		Schema: map[string]*schema.Schema{
			"users": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceMeIdentityUsersRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	users := make([]string, 0)

	err := config.OVHClient.Get(
		"/me/identity/user",
		&users,
	)
	if err != nil {
		return fmt.Errorf("Unable to get identity users:\n\t %q", err)
	}
	log.Printf("[DEBUG] identity users: %+v", users)

	sort.Strings(users)
	d.SetId(hashcode.Strings(users))
	d.Set("users", users)

	return nil
}

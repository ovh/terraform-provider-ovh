package ovh

import (
	"context"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceMeIdentityGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMeIdentityGroupsRead,
		Schema: map[string]*schema.Schema{
			"groups": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceMeIdentityGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	var groups []string

	err := config.OVHClient.GetWithContext(
		ctx,
		"/me/identity/group",
		&groups,
	)
	if err != nil {
		return diag.Errorf("Unable to get identity groups:\n\t %q", err)
	}
	log.Printf("[DEBUG] identity groups: %+v", groups)

	sort.Strings(groups)
	d.SetId(hashcode.Strings(groups))
	d.Set("groups", groups)

	return nil
}

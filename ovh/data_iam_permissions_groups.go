package ovh

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

func dataSourceIamPermissionsGroups() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"urns": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		ReadContext: datasourceIamGroupsRead,
	}
}

func datasourceIamGroupsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	var permissionsGroups []IamPermissionsGroup
	err := config.OVHClient.GetWithContext(ctx, "/v2/iam/permissionsGroup", &permissionsGroups)
	if err != nil {
		return diag.FromErr(err)
	}

	var permGrpUrns []string
	for _, p := range permissionsGroups {
		permGrpUrns = append(permGrpUrns, p.Urn)
	}

	d.Set("urns", permGrpUrns)

	sort.Strings(permGrpUrns)
	d.SetId(hashcode.Strings(permGrpUrns))
	return nil
}

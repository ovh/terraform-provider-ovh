package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceIamResourceGroups() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_groups": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
		ReadContext: datasourceIamResourceGroupsRead,
	}
}

func datasourceIamResourceGroupsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	var groups []IamResourceGroup
	err := config.OVHClient.GetWithContext(ctx, "/v2/iam/resourceGroup?details=true", &groups)
	if err != nil {
		return diag.FromErr(err)
	}

	var grpsId []string
	for _, grp := range groups {
		grpsId = append(grpsId, grp.ID)
	}

	d.SetId(hashcode.Strings(grpsId))
	d.Set("resource_groups", grpsId)
	return nil
}

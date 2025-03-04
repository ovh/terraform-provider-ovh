package ovh

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

func dataSourceIamPolicies() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"policies": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		ReadContext: datasourceIamPoliciesRead,
	}
}

func datasourceIamPoliciesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	var policies []IamPolicy
	err := config.OVHClient.GetWithContext(ctx, "/v2/iam/policy", &policies)
	if err != nil {
		return diag.FromErr(err)
	}

	var polIDs []string
	for _, p := range policies {
		polIDs = append(polIDs, p.Id)
	}

	d.Set("policies", polIDs)

	sort.Strings(polIDs)
	d.SetId(hashcode.Strings(polIDs))
	return nil
}

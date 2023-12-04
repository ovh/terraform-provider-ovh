package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIamResourceGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resources": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		ReadContext: datasourceIamResourceGroupRead,
	}
}

func datasourceIamResourceGroupRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	id := d.Get("id").(string)

	var pol IamResourceGroup
	err := config.OVHClient.GetWithContext(ctx, fmt.Sprintf("/v2/iam/resourceGroup/%s?details=true", url.PathEscape(id)), &pol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	var urns []string
	for _, r := range pol.Resources {
		urns = append(urns, r.URN)
	}

	d.Set("resources", urns)
	d.Set("name", pol.Name)
	d.Set("owner", pol.Owner)
	d.Set("created_at", pol.CreatedAt)
	d.Set("updated_at", pol.UpdatedAt)
	d.Set("read_only", pol.ReadOnly)
	d.Set("urn", pol.URN)

	return nil
}

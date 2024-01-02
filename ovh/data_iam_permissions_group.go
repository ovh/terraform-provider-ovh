package ovh

import (
	"context"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIamPermissionsGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"urn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"allow": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"except": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"deny": {
				Type:     schema.TypeSet,
				Optional: true,
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
				Optional: true,
				Computed: true,
			},
		},
		ReadContext: datasourceIamPermissionsGroupRead,
	}
}

func datasourceIamPermissionsGroupRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	id := d.Get("urn").(string)

	var pol IamPermissionsGroup
	err := config.OVHClient.GetWithContext(ctx, "/v2/iam/permissionsGroup/"+url.PathEscape(id), &pol)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range pol.ToMap() {
		err := d.Set(k, v)
		if err != nil {
			return diag.Errorf("key: %s; value: %v; err: %v", k, v, err)
		}
	}
	d.SetId(id)
	return nil
}

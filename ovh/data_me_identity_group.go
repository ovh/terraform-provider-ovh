package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMeIdentityGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMeIdentityGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_group": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"creation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMeIdentityGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	group := d.Get("name").(string)

	endpoint := fmt.Sprintf("/me/identity/group/%s", url.PathEscape(group))

	var resp MeIdentityGroupResponse
	err := config.OVHClient.GetWithContext(
		ctx,
		endpoint,
		&resp,
	)
	if err != nil {
		return diag.Errorf("Unable to get identity group detail:\n\t %q", err)
	}

	log.Printf("[DEBUG] identity groups: %+v", resp)

	d.SetId(group)

	d.Set("name", resp.Name)
	d.Set("default_group", resp.DefaultGroup)
	d.Set("last_update", resp.LastUpdate)
	d.Set("creation", resp.Creation)
	d.Set("description", resp.Description)
	d.Set("role", resp.Role)
	d.Set("urn", resp.URN)

	return nil
}

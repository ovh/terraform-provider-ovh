package ovh

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

func dataSourceMeApiOauth2Clients() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMeApiOauth2ClientsRead,
		Schema: map[string]*schema.Schema{
			"client_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceMeApiOauth2ClientsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	var oAuth2ClientsIds []string
	err := config.OVHClient.GetWithContext(ctx, "/me/api/oauth2/client", &oAuth2ClientsIds)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("client_ids", oAuth2ClientsIds)

	sort.Strings(oAuth2ClientsIds)
	d.SetId(hashcode.Strings(oAuth2ClientsIds))
	return nil
}

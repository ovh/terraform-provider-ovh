package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMeIdentityProvider() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMeIdentityProviderRead,

		Schema: map[string]*schema.Schema{
			"group_attribute_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"requested_attributes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_required": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name_format": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
						},
					},
				},
			},
			"idp_signing_certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expiration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subject": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"sso_service_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_attribute_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disable_users": {
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
		},
	}
}

func dataSourceMeIdentityProviderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	providerConfDetails := &MeIdentityProviderResponse{}
	if err := config.OVHClient.GetWithContext(ctx, "/me/identity/provider", providerConfDetails); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("ovh_sso")
	d.Set("group_attribute_name", providerConfDetails.GroupAttributeName)
	d.Set("disable_users", providerConfDetails.DisableUsers)
	d.Set("requested_attributes", requestedAttributesToMapList(providerConfDetails.Extensions.RequestedAttributes))
	d.Set("idp_signing_certificates", idpSigningCertificatesToMapList(providerConfDetails.IdpSigningCertificates))
	d.Set("sso_service_url", providerConfDetails.SsoServiceUrl)
	d.Set("user_attribute_name", providerConfDetails.UserAttributeName)
	d.Set("creation", providerConfDetails.Creation)
	d.Set("last_update", providerConfDetails.LastUpdate)

	return nil
}

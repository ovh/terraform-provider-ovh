package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func dataSourceMeApiOauth2Client() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMeApiOauth2ClientRead,
		Schema: map[string]*schema.Schema{
			"callback_urls": {
				Type:        schema.TypeList,
				Description: "Callback URLs of the applications using this oauth2 client. Required if using the AUTHORIZATION_CODE flow.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"client_id": {
				Type:        schema.TypeString,
				Description: "Client ID for the oauth2 client, generated during the resource creation.",
				Required:    true,
			},
			"client_secret": {
				Type:        schema.TypeString,
				Description: "Secret for the oauth2 client, generated during the oauth2 client creation. Can be specified in the data resource.",
				Optional:    true,
				Sensitive:   true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A description of your oauth2 client.",
				Computed:    true,
			},
			"identity": {
				Type:        schema.TypeString,
				Description: "URN that will allow you to associate this oauth2 client with an access policy.",
				Computed:    true,
			},
			"flow": {
				Type:        schema.TypeString,
				Description: "OAuth2 flow type implemented for this oauth2 client. Can be either AUTHORIZATION_CODE or CLIENT_CREDENTIALS",
				Computed:    true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMeApiOauth2ClientRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceAccount := &ApiOauth2ClientReadResponse{}

	// Query the oauth2 client using its client ID
	endpoint := fmt.Sprintf("/me/api/oauth2/client/%s", url.PathEscape(d.Get("client_id").(string)))
	if err := config.OVHClient.GetWithContext(ctx, endpoint, serviceAccount); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	log.Printf("[DEBUG] Read oauth2 client: %s", endpoint)

	// Populate the state with the response body parameters
	d.SetId(serviceAccount.ClientId)
	d.Set("callback_urls", serviceAccount.CallbackUrls)
	d.Set("client_id", serviceAccount.ClientId)
	d.Set("description", serviceAccount.Description)
	d.Set("flow", serviceAccount.Flow)
	d.Set("identity", serviceAccount.Identity)
	d.Set("name", serviceAccount.Name)

	return nil
}

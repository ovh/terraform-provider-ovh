package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceApiOauth2Client() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApiOauth2ClientCreate,
		ReadContext:   resourceApiOauth2ClientRead,
		UpdateContext: resourceApiOauth2ClientUpdate,
		DeleteContext: resourceApiOauth2ClientDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

				// Get the resource id from the terraform import command.
				// If it contains a pipe character, try to parse it as client_id|client_secret
				client_config_string := d.Id()
				if strings.Contains(client_config_string, "|") {

					log.Printf("[INFO] Importing oauth2 client formatted as 'client_id|client_secret'")

					client_config_params := strings.Split(client_config_string, "|")
					if len(client_config_params) != 2 {
						return nil, fmt.Errorf("Error importing api oauth2 client: Resource IDs with the pipe character should be formatted as 'client_id|client_secret': %s", d.Id())
					}
					d.SetId(client_config_params[0])
					d.Set("client_secret", client_config_params[1])
				}

				// Use the provided resource id as the client_id
				d.Set("client_id", d.Id())
				d.Set("discard_client_secret", false)
				return []*schema.ResourceData{d}, nil
			},
		},
		// When a previously created oauth2Client did not keep the client_secret value and now wants to,
		// a new client must be created, so we force a new resource to be created.
		CustomizeDiff: customdiff.ForceNewIfChange("discard_client_secret", func(ctx context.Context, old, new, meta any) bool {
			return old.(bool) && !new.(bool)
		}),
		Schema: map[string]*schema.Schema{
			"callback_urls": {
				Type:        schema.TypeList,
				Description: "Callback URLs of the applications using this oauth2 client. Required if using the AUTHORIZATION_CODE flow.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"client_id": {
				Type:        schema.TypeString,
				Description: "Client ID for the oauth2 client, generated during the resource creation.",
				Computed:    true,
			},
			"client_secret": {
				Type:        schema.TypeString,
				Description: "Secret for the oauth2 client, generated during the oauth2 client creation.",
				Computed:    true,
				Sensitive:   true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A description of your oauth2 client.",
				Required:    true,
			},
			"identity": {
				Type:        schema.TypeString,
				Description: "URN that will allow you to associate this oauth2 client with an access policy",
				Computed:    true,
			},
			"flow": {
				Type:         schema.TypeString,
				ValidateFunc: helpers.ValidateEnum([]string{"AUTHORIZATION_CODE", "CLIENT_CREDENTIALS"}),
				Description:  "OAuth2 flow type implemented for this oauth2 client. Can be either AUTHORIZATION_CODE or CLIENT_CREDENTIALS",
				ForceNew:     true,
				Required:     true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"discard_client_secret": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Do not keep the client_secret value in Terraform state. The client secret will not be available for use in Terraform",
			},
		},
	}
}

// Common function with the datasource
func resourceApiOauth2ClientRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceAccount := &ApiOauth2ClientReadResponse{}

	// Query the oauth2 client using its client ID
	endpoint := fmt.Sprintf("/me/api/oauth2/client/%s", url.PathEscape(d.Get("client_id").(string)))
	if err := config.OVHClient.GetWithContext(ctx, endpoint, serviceAccount); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	// Populate the state with the response body parameters
	d.Set("callback_urls", serviceAccount.CallbackUrls)
	d.Set("client_id", serviceAccount.ClientId)
	d.Set("description", serviceAccount.Description)
	d.Set("flow", serviceAccount.Flow)
	d.Set("identity", serviceAccount.Identity)
	d.Set("name", serviceAccount.Name)

	return nil
}

func resourceApiOauth2ClientCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Declare an empty array if no callback url is provided
	callbackUrls := []string{}
	if d.Get("callback_urls") != nil {
		itemsRaw := d.Get("callback_urls").([]interface{})
		for _, raw := range itemsRaw {
			callbackUrls = append(callbackUrls, raw.(string))
		}
	}

	params := &ApiOauth2ClientCreateOpts{
		CallbackUrls: callbackUrls,
		Description:  d.Get("description").(string),
		Flow:         d.Get("flow").(string),
		Name:         d.Get("name").(string),
	}

	log.Printf("[DEBUG] Will create api oauth2 client: %s", params.Name)

	// Create the oauth2 client using the provided parameters
	var response ApiOauth2ClientCreateResponse
	err := config.OVHClient.PostWithContext(ctx, "/me/api/oauth2/client", params, &response)
	if err != nil {
		return diag.Errorf("Error creating api oauth2 client %s:\n\t %q", params.Name, err)
	}

	// Use the client id as the resource id in the state
	d.SetId(response.ClientId)

	// Populate the state with the response body parameters
	d.Set("client_id", response.ClientId)
	// If discard_client_secret is true, we do not persist the client_secret in state. It will not be kept or usable in other steps.
	discardClientSecret := d.Get("discard_client_secret").(bool)
	if !discardClientSecret {
		d.Set("client_secret", response.ClientSecret)
	}

	return resourceApiOauth2ClientRead(ctx, d, meta)
}

func resourceApiOauth2ClientUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Update the client_secret based on the "discard_client_secret" param.
	oldDiscard, newDiscard := d.GetChange("discard_client_secret")
	// From false to true: remove the secret from state
	if d.HasChange("discard_client_secret") && !oldDiscard.(bool) && newDiscard.(bool) {
		d.Set("client_secret", nil) // remove the value from the state
		log.Printf("[DEBUG] client_secret value has been removed from terraform state")
	}

	// Only call the API when an attribute it knows about has changed
	if !d.HasChanges("callback_urls", "description", "name") {
		return resourceApiOauth2ClientRead(ctx, d, meta)
	}

	// Declare an empty array if no callback url is provided
	callbackUrls := []string{}
	if d.Get("callback_urls") != nil {
		itemsRaw := d.Get("callback_urls").([]interface{})
		for _, raw := range itemsRaw {
			callbackUrls = append(callbackUrls, raw.(string))
		}
	}

	params := &ApiOauth2ClientUpdateOpts{
		CallbackUrls: callbackUrls,
		Description:  d.Get("description").(string),
		Name:         d.Get("name").(string),
	}

	err := config.OVHClient.PutWithContext(ctx,
		fmt.Sprintf("/me/api/oauth2/client/%s", url.PathEscape(d.Id())),
		params,
		nil,
	)
	if err != nil {
		return diag.Errorf("Unable to update api oauth2 client %s:\n\t %q", url.PathEscape(d.Id()), err)
	}

	log.Printf("[DEBUG] Updated api oauth2 client %s", d.Id())
	return resourceApiOauth2ClientRead(ctx, d, meta)
}

func resourceApiOauth2ClientDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	err := config.OVHClient.DeleteWithContext(ctx,
		fmt.Sprintf("/me/api/oauth2/client/%s", url.PathEscape(d.Id())),
		nil,
	)
	if err != nil {
		return diag.Errorf("Unable to delete api oauth2 client %s:\n\t %q", d.Id(), err)
	}

	log.Printf("[DEBUG] Deleted api oauth2 client %s", d.Id())
	d.SetId("")
	return nil
}

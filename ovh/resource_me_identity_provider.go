package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMeIdentityProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMeIdentityProviderCreate,
		ReadContext:   resourceMeIdentityProviderRead,
		UpdateContext: resourceMeIdentityProviderUpdate,
		DeleteContext: resourceMeIdentityProviderDelete,

		// Importer is voluntarily disabled for this resource. As the `metadata`
		// attribute is not retrievable with GET /me/identity/provider, there
		// is no way to create a valid resource using `terraform import`
		Importer: nil,

		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_attribute_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"requested_attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_required": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name_format": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required: true,
						},
					},
				},
			},
			"disable_users": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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

func resourceMeIdentityProviderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	providerConfDetails := &MeIdentityProviderResponse{}
	if err := config.OVHClient.GetWithContext(ctx, "/me/identity/provider", providerConfDetails); err != nil {
		return diag.FromErr(err)
	}

	d.Set("group_attribute_name", providerConfDetails.GroupAttributeName)
	d.Set("disable_users", providerConfDetails.DisableUsers)
	d.Set("requested_attributes", requestedAttributesToMapList(providerConfDetails.Extensions.RequestedAttributes))
	d.Set("creation", providerConfDetails.Creation)
	d.Set("last_update", providerConfDetails.LastUpdate)

	return nil
}

func resourceMeIdentityProviderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	metadata := d.Get("metadata").(string)

	groupAttributeName := d.Get("group_attribute_name").(string)
	disableUsers := d.Get("disable_users").(bool)
	requestedAttributes, err := loadMeIdentityProviderAttributeListFromResource(d.Get("requested_attributes"))
	if err != nil {
		return diag.FromErr(err)
	}

	params := &MeIdentityProviderCreateOpts{
		Metadata:           metadata,
		GroupAttributeName: groupAttributeName,
		DisableUsers:       disableUsers,
		Extensions: MeIdentityProviderExtensions{
			RequestedAttributes: requestedAttributes,
		},
	}

	err = config.OVHClient.PostWithContext(ctx, "/me/identity/provider", params, nil)
	if err != nil {
		return diag.Errorf("Error creating identity provider:\n\t %v", err)
	}

	// As there is only one Identity Provider configurable, we use a constant ID
	d.SetId("ovh_sso")

	return resourceMeIdentityProviderRead(ctx, d, meta)
}

func resourceMeIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	groupAttributeName := d.Get("group_attribute_name").(string)
	disableUsers := d.Get("disable_users").(bool)
	requestedAttributes, err := loadMeIdentityProviderAttributeListFromResource(d.Get("requested_attributes"))
	if err != nil {
		return diag.FromErr(err)
	}

	params := &MeIdentityProviderUpdateOpts{
		GroupAttributeName: groupAttributeName,
		DisableUsers:       disableUsers,
		Extensions: MeIdentityProviderExtensions{
			RequestedAttributes: requestedAttributes,
		},
	}
	err = config.OVHClient.PutWithContext(ctx,
		"/me/identity/provider",
		params,
		nil,
	)
	if err != nil {
		return diag.Errorf("Unable to update identity provider:\n\t %q", err)
	}

	return resourceMeIdentityProviderRead(ctx, d, meta)
}

func resourceMeIdentityProviderDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	err := config.OVHClient.DeleteWithContext(ctx,
		"/me/identity/provider",
		nil,
	)
	if err != nil {
		return diag.Errorf("Unable to delete identity provider:\n\t %q", err)
	}

	d.SetId("")
	return nil
}

package ovh

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceMeIdentityGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMeIdentityGroupCreate,
		ReadContext:   resourceMeIdentityGroupRead,
		UpdateContext: resourceMeIdentityGroupUpdate,
		DeleteContext: resourceMeIdentityGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("name", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role": {
				Type:         schema.TypeString,
				ValidateFunc: helpers.ValidateEnum([]string{"ADMIN", "REGULAR", "UNPRIVILEGED", "NONE"}),
				Optional:     true,
				Default:      "NONE",
			},
			"default_group": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"last_update": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Common function with the datasource
func resourceMeIdentityGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	identityGroup := &MeIdentityGroupResponse{}

	endpoint := fmt.Sprintf("/me/identity/group/%s", d.Id())
	if err := config.OVHClient.GetWithContext(ctx, endpoint, identityGroup); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	d.Set("default_group", identityGroup.DefaultGroup)
	d.Set("last_update", identityGroup.LastUpdate)
	d.Set("creation", identityGroup.Creation)
	d.Set("description", identityGroup.Description)
	d.Set("role", identityGroup.Role)

	d.Set("urn", fmt.Sprintf("urn:v1:%s:identity:group:%s/%s", config.Plate, config.Account, identityGroup.Name))

	return nil
}

func resourceMeIdentityGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	role := d.Get("role").(string)

	params := &MeIdentityGroupCreateOpts{
		Name:        name,
		Description: description,
		Role:        role,
	}

	log.Printf("[DEBUG] Will create identity group: %s", params.Name)

	err := config.OVHClient.PostWithContext(ctx, "/me/identity/group", params, nil)
	if err != nil {
		return diag.Errorf("Error creating identity group %s:\n\t %q", params.Name, err)
	}

	d.SetId(name)

	return resourceMeIdentityGroupRead(ctx, d, meta)
}

func resourceMeIdentityGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	description := d.Get("description").(string)
	role := d.Get("role").(string)

	params := &MeIdentityGroupUpdateOpts{
		Description: description,
		Role:        role,
	}
	err := config.OVHClient.PutWithContext(ctx,
		fmt.Sprintf("/me/identity/group/%s", d.Id()),
		params,
		nil,
	)
	if err != nil {
		return diag.Errorf("Unable to update identity group %s:\n\t %q", d.Id(), err)
	}

	log.Printf("[DEBUG] Updated identity group %s", d.Id())
	return resourceMeIdentityGroupRead(ctx, d, meta)
}

func resourceMeIdentityGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	err := config.OVHClient.DeleteWithContext(ctx,
		fmt.Sprintf("/me/identity/group/%s", d.Id()),
		nil,
	)
	if err != nil {
		return diag.Errorf("Unable to delete identity group %s:\n\t %q", d.Id(), err)
	}

	log.Printf("[DEBUG] Deleted identity group %s", d.Id())
	d.SetId("")
	return nil
}

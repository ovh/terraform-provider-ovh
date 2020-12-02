package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceMeIdentityUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceMeIdentityUserCreate,
		Read:   resourceMeIdentityUserRead,
		Update: resourceMeIdentityUserUpdate,
		Delete: resourceMeIdentityUserDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User description",
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User's email",
			},
			"group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User's group",
			},
			"login": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User's login suffix",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "User's password",
			},
			"creation": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date of this user",
			},
			"last_update": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update of this user",
			},
			"password_last_update": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the user changed his password for the last time",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current user's status",
			},
		},
	}
}

// Common function with the datasource
func resourceMeIdentityUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	identityUser := &MeIdentityUserResponse{}

	endpoint := fmt.Sprintf("/me/identity/user/%s", d.Id())
	if err := config.OVHClient.Get(endpoint, identityUser); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("login", identityUser.Login)
	d.Set("creation", identityUser.Creation)
	d.Set("description", identityUser.Description)
	d.Set("email", identityUser.Email)
	d.Set("group", identityUser.Group)
	d.Set("last_update", identityUser.LastUpdate)
	d.Set("password_last_update", identityUser.PasswordLastUpdate)
	d.Set("status", identityUser.Status)

	return nil
}

func resourceMeIdentityUserCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	description := d.Get("description").(string)
	email := d.Get("email").(string)
	group := d.Get("group").(string)
	login := d.Get("login").(string)
	password := d.Get("password").(string)
	params := &MeIdentityUserCreateOpts{
		Description: description,
		Email:       email,
		Group:       group,
		Login:       login,
		Password:    password,
	}

	log.Printf("[DEBUG] Will create identity user: %s", params.Email)

	err := config.OVHClient.Post("/me/identity/user", params, nil)
	if err != nil {
		return fmt.Errorf("Error creating identity user %s:\n\t %q", params.Email, err)
	}

	d.SetId(login)

	return resourceMeIdentityUserRead(d, meta)
}

func resourceMeIdentityUserUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id := d.Id()
	description := d.Get("description").(string)
	email := d.Get("email").(string)
	group := d.Get("group").(string)
	login := d.Get("login").(string)
	params := &MeIdentityUserUpdateOpts{
		Login:       login,
		Description: description,
		Email:       email,
		Group:       group,
	}
	err := config.OVHClient.Put(
		fmt.Sprintf("/me/identity/user/%s", id),
		nil,
		params,
	)
	if err != nil {
		return fmt.Errorf("Unable to update identity user %s:\n\t %q", id, err)
	}

	log.Printf("[DEBUG] Updated identity user %s", id)
	return resourceMeIdentityUserRead(d, meta)
}

func resourceMeIdentityUserDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id := d.Id()
	err := config.OVHClient.Delete(
		fmt.Sprintf("/me/identity/user/%s", id),
		nil,
	)
	if err != nil {
		return fmt.Errorf("Unable to delete identity user %s:\n\t %q", id, err)
	}

	log.Printf("[DEBUG] Deleted identity user %s", id)
	d.SetId("")
	return nil
}

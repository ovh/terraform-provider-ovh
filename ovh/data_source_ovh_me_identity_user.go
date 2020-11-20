package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMeIdentityUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMeIdentityUserRead,
		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User's login",
			},
			"login": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User's login suffix",
			},
			"creation": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date of this user",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User description",
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User's email",
			},
			"group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User's group",
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

func dataSourceMeIdentityUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	identityUser := &MeIdentityUserResponse{}

	user := d.Get("user").(string)
	err := config.OVHClient.Get(
		fmt.Sprintf("/me/identity/user/%s", user),
		identityUser,
	)
	if err != nil {
		return fmt.Errorf("Unable to find identity user %s:\n\t %q", user, err)
	}
	log.Printf("[DEBUG] identity user for %s: %+v", user, identityUser)

	d.SetId(user)
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

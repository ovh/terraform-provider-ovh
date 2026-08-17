package ovh

import (
	"fmt"
	"log"
	"net/url"

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
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
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
			"groups": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Additional groups the user belongs to (other than the main group)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
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
		return fmt.Errorf("unable to find identity user %s:\n\t %q", user, err)
	}

	d.Set("urn", fmt.Sprintf("urn:v1:%s:identity:user:%s/%s", config.Plate, config.Account, user))

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

	// Discover additional group memberships by listing all groups
	// and checking which ones contain this user
	var allGroups []string
	if err := config.OVHClient.Get("/me/identity/group", &allGroups); err != nil {
		return fmt.Errorf("unable to list identity groups:\n\t %q", err)
	}

	var memberGroups []string
	for _, groupName := range allGroups {
		// Skip the user's main group
		if groupName == identityUser.Group {
			continue
		}
		var users []string
		groupEndpoint := fmt.Sprintf("/me/identity/group/%s/user", url.PathEscape(groupName))
		if err := config.OVHClient.Get(groupEndpoint, &users); err != nil {
			log.Printf("[WARN] Could not read users for group %s: %s", groupName, err)
			continue
		}
		for _, u := range users {
			if u == identityUser.Login {
				memberGroups = append(memberGroups, groupName)
				break
			}
		}
	}
	d.Set("groups", memberGroups)

	return nil
}

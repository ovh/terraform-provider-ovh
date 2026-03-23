package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
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
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User's email",
			},
			"group": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "DEFAULT",
				Description: "User's main group",
			},
			"groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Additional groups the user belongs to (other than the main group)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
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

	d.Set("urn", fmt.Sprintf("urn:v1:%s:identity:user:%s/%s", config.Plate, config.Account, identityUser.Login))

	d.Set("login", identityUser.Login)
	d.Set("creation", identityUser.Creation)
	d.Set("description", identityUser.Description)
	d.Set("email", identityUser.Email)
	d.Set("group", identityUser.Group)
	d.Set("last_update", identityUser.LastUpdate)
	d.Set("password_last_update", identityUser.PasswordLastUpdate)
	d.Set("status", identityUser.Status)

	// Discover all additional group memberships by listing all groups
	// and checking which ones contain this user
	var allGroups []string
	if err := config.OVHClient.Get("/me/identity/group", &allGroups); err != nil {
		log.Printf("[WARN] Could not list identity groups: %s", err)
	} else {
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
	}

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

	// Add user to additional groups
	if v, ok := d.GetOk("groups"); ok {
		for _, g := range v.(*schema.Set).List() {
			groupName := g.(string)
			if err := addUserToGroup(config, groupName, login); err != nil {
				return err
			}
		}
	}

	return resourceMeIdentityUserRead(d, meta)
}

func resourceMeIdentityUserUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id := d.Id()
	description := d.Get("description").(string)
	email := d.Get("email").(string)
	group := d.Get("group").(string)
	params := &MeIdentityUserUpdateOpts{
		Description: description,
		Email:       email,
		Group:       group,
	}
	err := config.OVHClient.Put(
		fmt.Sprintf("/me/identity/user/%s", id),
		params,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Unable to update identity user %s:\n\t %q", id, err)
	}

	log.Printf("[DEBUG] Updated identity user %s", id)

	if d.HasChange("groups") {
		old, new := d.GetChange("groups")
		oldSet := old.(*schema.Set)
		newSet := new.(*schema.Set)

		toAdd := newSet.Difference(oldSet)
		toRemove := oldSet.Difference(newSet)

		for _, g := range toAdd.List() {
			if err := addUserToGroup(config, g.(string), id); err != nil {
				return err
			}
		}
		for _, g := range toRemove.List() {
			if err := removeUserFromGroup(config, g.(string), id); err != nil {
				return err
			}
		}
	}

	return resourceMeIdentityUserRead(d, meta)
}

func resourceMeIdentityUserDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id := d.Id()

	// Remove user from all additional groups before deleting
	if v, ok := d.GetOk("groups"); ok {
		for _, g := range v.(*schema.Set).List() {
			if err := removeUserFromGroup(config, g.(string), id); err != nil {
				log.Printf("[WARN] Could not remove user %s from group %s: %s", id, g.(string), err)
			}
		}
	}

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

func addUserToGroup(config *Config, group, login string) error {
	params := map[string]string{"user": login}
	endpoint := fmt.Sprintf("/me/identity/group/%s/user", url.PathEscape(group))
	if err := config.OVHClient.Post(endpoint, params, nil); err != nil {
		return fmt.Errorf("Error adding user %s to group %s:\n\t %q", login, group, err)
	}
	log.Printf("[DEBUG] Added user %s to group %s", login, group)
	return nil
}

func removeUserFromGroup(config *Config, group, login string) error {
	endpoint := fmt.Sprintf("/me/identity/group/%s/user/%s", url.PathEscape(group), url.PathEscape(login))
	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return fmt.Errorf("Error removing user %s from group %s:\n\t %q", login, group, err)
	}
	log.Printf("[DEBUG] Removed user %s from group %s", login, group)
	return nil
}

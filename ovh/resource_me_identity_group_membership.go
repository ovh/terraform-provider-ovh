package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceMeIdentityGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceMeIdentityGroupMembershipCreate,
		Read:   resourceMeIdentityGroupMembershipRead,
		Delete: resourceMeIdentityGroupMembershipDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMeIdentityGroupMembershipImportState,
		},

		Schema: map[string]*schema.Schema{
			"login": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The login of the identity user to add to the group.",
			},
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the identity group to add the user to.",
			},
		},
	}
}

func resourceMeIdentityGroupMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	login := d.Get("login").(string)
	group := d.Get("group").(string)

	if err := addUserToGroup(config, group, login); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", login, group))
	return resourceMeIdentityGroupMembershipRead(d, meta)
}

func resourceMeIdentityGroupMembershipRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	login := d.Get("login").(string)
	group := d.Get("group").(string)

	var users []string
	endpoint := fmt.Sprintf("/me/identity/group/%s/user", url.PathEscape(group))
	if err := config.OVHClient.Get(endpoint, &users); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	found := false
	for _, u := range users {
		if u == login {
			found = true
			break
		}
	}

	if !found {
		d.SetId("")
		return nil
	}

	d.Set("login", login)
	d.Set("group", group)

	return nil
}

func resourceMeIdentityGroupMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	login := d.Get("login").(string)
	group := d.Get("group").(string)

	if err := removeUserFromGroup(config, group, login); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceMeIdentityGroupMembershipImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not login/group formatted")
	}
	login := splitId[0]
	group := splitId[1]

	d.SetId(fmt.Sprintf("%s/%s", login, group))
	d.Set("login", login)
	d.Set("group", group)

	return []*schema.ResourceData{d}, nil
}

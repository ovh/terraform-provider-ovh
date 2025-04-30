package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"golang.org/x/exp/slices"

	"github.com/ovh/go-ovh/ovh"
)

func resourceCloudProjectUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectUserCreate,
		Read:   resourceCloudProjectUserRead,
		Delete: resourceCloudProjectUserDelete,
		Update: resourceCloudProjectUserUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCloudProjectUserImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
				Description: "Service name of the resource representing the id of the cloud project.",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"role_names": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"rotate_when_changed": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    false,
				Description: "A map of arbitrary key/value pairs that will trigger password regeneration when they change, enabling password rotation based on external conditions such as a rotating timestamp.",
			},

			// Computed
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"openstack_rc": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permissions": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func validateCloudProjectUserRoleFunc(config *Config, serviceName string, roles []string, role string) (*CloudProjectrolesResponse, error) {
	endpoint := fmt.Sprintf("/cloud/project/%s/role",
		url.PathEscape(serviceName),
	)
	res := &CloudProjectrolesResponse{}
	if err := config.OVHClient.Get(endpoint, res); err != nil {
		return nil, fmt.Errorf("calling Get %s", endpoint)
	}

	ovhRole := make([]string, 0, len(res.Roles))
	for _, val := range res.Roles {
		ovhRole = append(ovhRole, val.Name)
	}

	rolesToCheck := roles
	if role != "" {
		rolesToCheck = append(rolesToCheck, role)
	}
	for _, role := range rolesToCheck {
		if !slices.Contains(ovhRole, role) {
			return nil, fmt.Errorf("role %q does not exist", role)
		}
	}

	return res, nil
}

func resourceCloudProjectUserImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	userId := d.Id()
	// Fallback to the environment variable if service_name not given
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE")

	splitId := strings.SplitN(userId, "/", 2)
	if len(splitId) == 2 {
		serviceName = splitId[0]
		userId = splitId[1]
	}

	if serviceName == "" || userId == "" {
		return nil, fmt.Errorf("Import Id is not service_name/id formatted")
	}

	d.SetId(userId)
	d.Set("service_name", serviceName)

	return []*schema.ResourceData{d}, nil
}

func resourceCloudProjectUserUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	userId := d.Id()

	// Check if rotate_when_changed has been modified
	if d.HasChange("rotate_when_changed") {
		log.Printf("[DEBUG] rotate_when_changed has been modified, regenerating password for user %s", userId)

		// Call the regeneratePassword endpoint
		regenerate, err := cloudProjectUserRegeneratePassword(config, serviceName, userId)
		if err != nil {
			return fmt.Errorf("error regenerating password for user %s: %s", userId, err)
		}

		// Update the password in the state
		if err := d.Set("password", regenerate.Password); err != nil {
			return fmt.Errorf("error setting password in state: %s", err)
		}

		log.Printf("[DEBUG] Successfully updated password in state for user %s", userId)
	}

	// Update roles if they changed
	if d.HasChange("role_name") || d.HasChange("role_names") {
		role := d.Get("role_name")
		roles, err := helpers.StringsFromSchema(d, "role_names")
		res := &CloudProjectrolesResponse{}

		res, err = validateCloudProjectUserRoleFunc(config, serviceName, roles, role.(string))
		if err != nil {
			return err
		}
		update := []string{}
		for _, i := range res.Roles {
			if slices.Contains(roles, i.Name) {
				update = append(update, i.Id)
			}
			if role == i.Name && !slices.Contains(update, i.Name) {
				update = append(update, i.Id)
			}
		}

		log.Printf("[DEBUG] roles IDs %s", update)
		log.Printf("[DEBUG] user %s", userId)
		endpoint := fmt.Sprintf("/cloud/project/%s/user/%s/role",
			url.PathEscape(serviceName),
			url.PathEscape(userId),
		)

		r := &CloudProjectUser{}
		data := &CloudProjectroleUpdate{
			RolesIds: update,
		}
		if err := config.OVHClient.Put(endpoint, data, r); err != nil {
			return fmt.Errorf("calling %s with params %s:\n\t %q", endpoint, data, err)
		}
	}

	return resourceCloudProjectUserRead(d, meta)
}

func resourceCloudProjectUserCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	role := d.Get("role_name")
	params := (&CloudProjectUserCreateOpts{}).FromResource(d)

	r := &CloudProjectUser{}
	_, err := validateCloudProjectUserRoleFunc(config, serviceName, params.Roles, role.(string))
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Will create public cloud user: %s", params)
	endpoint := fmt.Sprintf(
		"/cloud/project/%s/user",
		url.PathEscape(serviceName),
	)

	if err := config.OVHClient.Post(endpoint, params, r); err != nil {
		return fmt.Errorf("calling Post %s with params %s:\n\t %q", endpoint, params, err)
	}

	// Set Password only at creation time
	d.Set("password", r.Password)
	d.SetId(strconv.Itoa(r.Id))

	log.Printf("[DEBUG] Waiting for User %s:", r)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"ok"},
		Refresh:    waitForCloudProjectUser(config.OVHClient, serviceName, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for user (%s): %s", params, err)
	}
	log.Printf("[DEBUG] Created User %s", r)

	return resourceCloudProjectUserRead(d, meta)
}

func resourceCloudProjectUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	user := &CloudProjectUser{}

	log.Printf("[DEBUG] Will read public cloud user %s from project: %s", d.Id(), serviceName)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/user/%s",
		url.PathEscape(serviceName),
		d.Id(),
	)

	if err := config.OVHClient.Get(endpoint, user); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId(strconv.Itoa(user.Id))
	// set resource attributes
	for k, v := range user.ToMap() {
		d.Set(k, v)
	}

	openstackrc := make(map[string]string)
	if err := cloudUserGetOpenstackRC(serviceName, d.Id(), config.OVHClient, openstackrc); err != nil {
		return fmt.Errorf("Reading openstack creds for user %s: %s", d.Id(), err)
	}

	d.Set("openstack_rc", &openstackrc)
	d.Set("service_name", serviceName)

	log.Printf("[DEBUG] Read Public Cloud User %s", user)
	return nil
}

func resourceCloudProjectUserDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	id := d.Id()

	log.Printf("[DEBUG] Will delete public cloud user %s from project: %s", id, serviceName)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/user/%s",
		url.PathEscape(serviceName),
		id,
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return fmt.Errorf("calling Delete %s:\n\t %q", endpoint, err)
	}

	log.Printf("[DEBUG] Deleting Public Cloud User %s from project %s:", id, serviceName)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting"},
		Target:     []string{"deleted"},
		Refresh:    waitForCloudProjectUser(config.OVHClient, serviceName, id),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Deleting Public Cloud user %s from project %s", id, serviceName)
	}
	log.Printf("[DEBUG] Deleted Public Cloud User %s from project %s", id, serviceName)

	d.SetId("")

	return nil
}

// cloudProjectUserRegeneratePassword regenerates the password for a cloud project user
func cloudProjectUserRegeneratePassword(config *Config, serviceName, userId string) (*CloudProjectUserRegenerate, error) {
	log.Printf("[DEBUG] Will regenerate password for user %s from project: %s", userId, serviceName)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/user/%s/regeneratePassword",
		url.PathEscape(serviceName),
		url.PathEscape(userId),
	)

	r := &CloudProjectUserRegenerate{}
	if err := config.OVHClient.Post(endpoint, nil, r); err != nil {
		return nil, fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
	}

	log.Printf("[DEBUG] Successfully regenerated password for user %s", userId)
	return r, nil
}

var cloudUserOSTenantName = regexp.MustCompile("export OS_TENANT_NAME=\"?([[:alnum:]]+)\"?")
var cloudUserOSTenantId = regexp.MustCompile("export OS_TENANT_ID=\"??([[:alnum:]]+)\"??")
var cloudUserOSAuthURL = regexp.MustCompile("export OS_AUTH_URL=\"??([[:^space:]]+)\"??")
var cloudUserOSUsername = regexp.MustCompile("export OS_USERNAME=\"?([[:alnum:]]+)\"?")

func cloudUserGetOpenstackRC(serviceName, id string, c *ovh.Client, rc map[string]string) error {
	log.Printf("[DEBUG] Will read public cloud user openstack rc for project: %s, id: %s", serviceName, id)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/user/%s/openrc?region=to_be_overriden",
		url.PathEscape(serviceName),
		id,
	)

	r := &CloudProjectUserOpenstackRC{}

	if err := c.Get(endpoint, r); err != nil {
		return fmt.Errorf("calling Get %s:\n\t %q", endpoint, err)
	}

	authURL := cloudUserOSAuthURL.FindStringSubmatch(r.Content)
	if authURL == nil {
		return fmt.Errorf("couln't extract OS_AUTH_URL from content: \n\t%s", r.Content)
	}
	tenantName := cloudUserOSTenantName.FindStringSubmatch(r.Content)
	if tenantName == nil {
		return fmt.Errorf("couln't extract OS_TENANT_NAME from content: \n\t%s", r.Content)
	}
	tenantId := cloudUserOSTenantId.FindStringSubmatch(r.Content)
	if tenantId == nil {
		return fmt.Errorf("couln't extract OS_TENANT_ID from content: \n\t%s", r.Content)
	}
	username := cloudUserOSUsername.FindStringSubmatch(r.Content)
	if username == nil {
		return fmt.Errorf("couln't extract OS_USERNAME from content: \n\t%s", r.Content)
	}

	rc["OS_AUTH_URL"] = authURL[1]
	rc["OS_TENANT_ID"] = tenantId[1]
	rc["OS_TENANT_NAME"] = tenantName[1]
	rc["OS_USERNAME"] = username[1]

	return nil
}

func waitForCloudProjectUser(c *ovh.Client, serviceName, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &CloudProjectUser{}
		endpoint := fmt.Sprintf(
			"/cloud/project/%s/user/%s",
			url.PathEscape(serviceName),
			id,
		)
		if err := c.Get(endpoint, r); err != nil {
			if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
				log.Printf("[DEBUG] user id %s on project %s deleted", id, serviceName)
				return r, "deleted", nil
			} else {
				return r, "", err
			}
		}

		log.Printf("[DEBUG] Pending User: %s", r)
		return r, r.Status, nil
	}
}

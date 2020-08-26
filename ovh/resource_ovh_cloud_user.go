package ovh

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-ovh/ovh/helpers"

	"github.com/ovh/go-ovh/ovh"
)

func resourceCloudUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudUserCreate,
		Read:   resourceCloudUserRead,
		Delete: resourceCloudUserDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				DefaultFunc:   schema.EnvDefaultFunc("OVH_PROJECT_ID", nil),
				Description:   "Id of the cloud project. DEPRECATED, use `service_name` instead",
				ConflictsWith: []string{"service_name"},
			},
			"service_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "Service name of the resource representing the id of the cloud project.",
				ConflictsWith: []string{"project_id"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"role_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateCloudUserRoleFunc,
			},
			"role_names": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func validateCloudUserRoleFunc(v interface{}, k string) (ws []string, errors []error) {
	err := helpers.ValidateStringEnum(v.(string), []string{
		"administrator",
		"ai_training_operator",
		"authentication",
		"backup_operator",
		"compute_operator",
		"image_operator",
		"infrastructure_supervisor",
		"network_operator",
		"network_security_operator",
		"objectstore_operator",
		"volume_operator",
	})

	if err != nil {
		errors = append(errors, err)
	}
	return
}

func resourceCloudUserCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectId := helpers.GetNilStringPointerFromData(d, "project_id")
	serviceNamePtr := helpers.GetNilStringPointerFromData(d, "service_name")

	if serviceNamePtr == nil && projectId != nil && *projectId != "" {
		serviceNamePtr = projectId
	}

	if serviceNamePtr == nil || *serviceNamePtr == "" {
		return fmt.Errorf("service_name attribute is mandatory.")
	}

	serviceName := *serviceNamePtr

	params := (&CloudUserCreateOpts{}).FromResource(d)

	for _, role := range params.Roles {
		if _, errs := validateCloudUserRoleFunc(role, ""); errs != nil {
			return fmt.Errorf("roles contains unsupported value: %s.", role)
		}
	}

	r := &CloudUser{}

	log.Printf("[DEBUG] Will create public cloud user: %s", params)
	endpoint := fmt.Sprintf(
		"/cloud/project/%s/user",
		url.PathEscape(serviceName),
	)
	err := config.OVHClient.Post(endpoint, params, r)
	if err != nil {
		return fmt.Errorf("calling Post %s with params %s:\n\t %q", endpoint, params, err)
	}

	// Set Password only at creation time
	d.Set("password", r.Password)
	d.SetId(strconv.Itoa(r.Id))

	log.Printf("[DEBUG] Waiting for User %s:", r)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"ok"},
		Refresh:    waitForCloudUser(config.OVHClient, serviceName, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("waiting for user (%s): %s", params, err)
	}
	log.Printf("[DEBUG] Created User %s", r)

	return resourceCloudUserRead(d, meta)
}

func resourceCloudUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectId := helpers.GetNilStringPointerFromData(d, "project_id")
	serviceNamePtr := helpers.GetNilStringPointerFromData(d, "service_name")

	if serviceNamePtr == nil && projectId != nil && *projectId != "" {
		serviceNamePtr = projectId
	}

	if serviceNamePtr == nil || *serviceNamePtr == "" {
		return fmt.Errorf("service_name attribute is mandatory.")
	}

	serviceName := *serviceNamePtr

	user := &CloudUser{}

	log.Printf("[DEBUG] Will read public cloud user %s from project: %s", d.Id(), serviceName)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/user/%s",
		url.PathEscape(serviceName),
		d.Id(),
	)

	err := config.OVHClient.Get(endpoint, user)
	if err != nil {
		return fmt.Errorf("calling Get %s:\n\t %q", endpoint, err)
	}

	d.SetId(strconv.Itoa(user.Id))
	// set resource attributes
	for k, v := range user.ToMap() {
		d.Set(k, v)
	}

	openstackrc := make(map[string]string)
	err = cloudUserGetOpenstackRC(serviceName, d.Id(), config.OVHClient, openstackrc)
	if err != nil {
		return fmt.Errorf("Reading openstack creds for user %s: %s", d.Id(), err)
	}

	d.Set("openstack_rc", &openstackrc)
	log.Printf("[DEBUG] Read Public Cloud User %s", user)
	return nil
}

func resourceCloudUserDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectId := helpers.GetNilStringPointerFromData(d, "project_id")
	serviceNamePtr := helpers.GetNilStringPointerFromData(d, "service_name")

	if serviceNamePtr == nil && projectId != nil && *projectId != "" {
		serviceNamePtr = projectId
	}

	if serviceNamePtr == nil || *serviceNamePtr == "" {
		return fmt.Errorf("service_name attribute is mandatory.")
	}

	serviceName := *serviceNamePtr

	id := d.Id()

	log.Printf("[DEBUG] Will delete public cloud user %s from project: %s", id, serviceName)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/user/%s",
		url.PathEscape(serviceName),
		id,
	)

	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return fmt.Errorf("calling Delete %s:\n\t %q", endpoint, err)
	}

	log.Printf("[DEBUG] Deleting Public Cloud User %s from project %s:", id, serviceName)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting"},
		Target:     []string{"deleted"},
		Refresh:    waitForCloudUser(config.OVHClient, serviceName, id),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Deleting Public Cloud user %s from project %s", id, serviceName)
	}
	log.Printf("[DEBUG] Deleted Public Cloud User %s from project %s", id, serviceName)

	d.SetId("")

	return nil
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

	r := &CloudUserOpenstackRC{}

	err := c.Get(endpoint, r)
	if err != nil {
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

func waitForCloudUser(c *ovh.Client, serviceName, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &CloudUser{}
		endpoint := fmt.Sprintf(
			"/cloud/project/%s/user/%s",
			url.PathEscape(serviceName),
			id,
		)
		err := c.Get(endpoint, r)
		if err != nil {
			if err.(*ovh.APIError).Code == 404 {
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

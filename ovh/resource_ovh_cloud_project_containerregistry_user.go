package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceCloudProjectContainerRegistryUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectContainerRegistryUserCreate,
		Read:   resourceCloudProjectContainerRegistryUserRead,
		Delete: resourceCloudProjectContainerRegistryUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectContainerRegistryUserImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "Service name",
				ForceNew:    true,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"registry_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "RegistryID",
				Required:    true,
			},
			"login": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Registry name",
				Required:    true,
			},
			"email": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "User email.",
				Required:    true,
			},

			// Computed
			"password": {
				Type:        schema.TypeString,
				Description: "User password",
				Sensitive:   true,
				Computed:    true,
			},
			"user": {
				Type:        schema.TypeString,
				Description: "User name",
				Computed:    true,
			},
		},
	}
}

func resourceCloudProjectContainerRegistryUserImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 3)

	if len(splitId) != 3 {
		return nil, fmt.Errorf("Import Id is not service_name/registry_id/id id formatted")
	}

	serviceName := splitId[0]
	regId := splitId[1]
	id := splitId[2]

	d.SetId(id)
	d.Set("registry_id", regId)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectContainerRegistryUserCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	regId := d.Get("registry_id").(string)

	log.Printf("[DEBUG] Will create cloud project registry user for project: %s", serviceName)

	opts := (&CloudProjectContainerRegistryUserCreateOpts{}).FromResource(d)
	user := &CloudProjectContainerRegistryUser{}

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s/users",
		url.PathEscape(serviceName),
		url.PathEscape(regId),
	)

	err := config.OVHClient.Post(endpoint, opts, user)
	if err != nil {
		return fmt.Errorf("Error calling post %s:\n\t %q", endpoint, err)
	}

	for k, v := range user.ToMapWithKeys([]string{"email", "password", "user"}) {
		d.Set(k, v)
	}

	d.SetId(user.Id)

	return resourceCloudProjectContainerRegistryUserRead(d, meta)
}

func resourceCloudProjectContainerRegistryUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	regId := d.Get("registry_id").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will read cloud project registry users %s for project: %s", id, serviceName)

	users := []CloudProjectContainerRegistryUser{}

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s/users",
		url.PathEscape(serviceName),
		url.PathEscape(regId),
	)

	if err := config.OVHClient.Get(endpoint, &users); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	found := false

	for _, user := range users {
		if user.Id == id {
			found = true
			for k, v := range user.ToMapWithKeys([]string{"email", "user"}) {
				d.Set(k, v)
			}
		}
	}

	if !found {
		log.Printf("[DEBUG] cloud project registry user %s/%s for project %s not found", regId, id, serviceName)
		d.SetId("")
		return nil
	}

	return nil
}

func resourceCloudProjectContainerRegistryUserDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	regId := d.Get("registry_id").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will delete cloud project registry user %s/%s for project: %s", regId, id, serviceName)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s/users/%s",
		url.PathEscape(serviceName),
		url.PathEscape(regId),
		url.PathEscape(id),
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId("")

	return nil
}

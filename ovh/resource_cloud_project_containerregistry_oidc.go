package ovh

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

func resourceCloudProjectContainerRegistryOIDC() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectContainerRegistryOIDCCreate,
		Read:   resourceCloudProjectContainerRegistryOIDCRead,
		Delete: resourceCloudProjectContainerRegistryOIDCDelete,
		Update: resourceCloudProjectContainerRegistryOIDCUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectContainerRegistryOIDCImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(10 * time.Minute),
			Update:  schema.DefaultTimeout(10 * time.Minute),
			Delete:  schema.DefaultTimeout(10 * time.Minute),
			Read:    schema.DefaultTimeout(5 * time.Minute),
			Default: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"registry_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delete_users": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				ForceNew: true,
			},
			"oidc_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"oidc_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"oidc_client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"oidc_client_secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"oidc_scope": {
				Type:     schema.TypeString,
				Required: true,
			},
			"oidc_groups_claim": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"oidc_admin_group": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"oidc_verify_cert": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"oidc_auto_onboard": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
			},
			"oidc_user_claim": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
		},
	}
}

func resourceCloudProjectContainerRegistryOIDCImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	log.Printf("[DEBUG] Importing cloud project registry OIDC %s", givenId)

	splitId := strings.SplitN(givenId, "/", 3)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("import Id is not service_name/registryid formatted")
	}
	serviceName := splitId[0]
	registryID := splitId[1]
	d.SetId(serviceName + "/" + registryID)
	d.Set("registry_id", registryID)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectContainerRegistryOIDCCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	registryID := d.Get("registry_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/containerRegistry/%s/openIdConnect", serviceName, registryID)
	params := (&CloudProjectContainerRegistryOIDCCreateOpts{}).FromResource(d)
	res := &CloudProjectKubeOIDCResponse{}

	log.Printf("[DEBUG] Will create registry %s OIDC: %+v", registryID, params)
	err := config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		return fmt.Errorf("calling Post %s with params %v:\n\t %w", endpoint, params, err)
	}

	d.SetId(serviceName + "/" + registryID)

	log.Printf("[DEBUG] Registry %s OIDC created", registryID)

	return resourceCloudProjectContainerRegistryOIDCRead(d, meta)
}

func resourceCloudProjectContainerRegistryOIDCRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	registryID := d.Get("registry_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/containerRegistry/%s/openIdConnect", serviceName, registryID)
	res := &CloudProjectContainerRegistryOIDCResponse{}

	log.Printf("[DEBUG] Will read oidc from registry %s and project: %s", registryID, serviceName)
	err := config.OVHClient.Get(endpoint, res)
	if err != nil {
		if ovhErr, ok := err.(*ovh.APIError); ok && ovhErr.Code == 404 {
			// If the resource does not exist, remove it from the state to force recreation
			log.Printf("[DEBUG] Registry %s OIDC does not exist, removing from state", registryID)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("calling get %s %w", endpoint, err)
	}
	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(serviceName + "/" + registryID)
		}
	}

	log.Printf("[DEBUG] Read registry %+v", res)
	return nil
}

func resourceCloudProjectContainerRegistryOIDCUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	registryID := d.Get("registry_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/containerRegistry/%s/openIdConnect", serviceName, registryID)
	params := (&CloudProjectContainerRegistryOIDCUpdateOpts{}).FromResource(d)
	res := &CloudProjectContainerRegistryOIDCResponse{}

	log.Printf("[DEBUG] Will update registry %s OIDC: %+v", registryID, params)
	err := config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		return fmt.Errorf("calling Post %s with params %v:\n\t %w", endpoint, params, err)
	}

	log.Printf("[DEBUG] Registry %s OIDC updated", registryID)

	return nil
}

func resourceCloudProjectContainerRegistryOIDCDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	registryID := d.Get("registry_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/containerRegistry/%s/openIdConnect", serviceName, registryID)

	log.Printf("[DEBUG] Will delete registry %s OIDC", registryID)
	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return fmt.Errorf("calling delete %s %w", endpoint, err)
	}

	log.Printf("[DEBUG] Registry %s OIDC deleted", registryID)

	return nil
}

package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudProjectContainerRegistryIAM() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectContainerRegistryIAMCreate,
		Delete: resourceCloudProjectContainerRegistryIAMDelete,
		Read:   resourceCloudProjectContainerRegistryIAMRead,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectContainerRegistryIAMImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(10 * time.Minute),
			Delete:  schema.DefaultTimeout(10 * time.Minute),
			Read:    schema.DefaultTimeout(10 * time.Minute),
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
			"iam_enabled": {
				Type:     schema.TypeBool,
				Required: false,
				Computed: true,
			},
		},
	}
}

func resourceCloudProjectContainerRegistryIAMImportState(d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	givenId := d.Id()

	log.Printf("[DEBUG] Importing cloud project registry IAM %s", givenId)

	splitId := strings.SplitN(givenId, "/", 3)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("import Id is not service_name/registryId formatted")
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

func resourceCloudProjectContainerRegistryIAMRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	registryID := d.Get("registry_id").(string)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s",
		url.PathEscape(serviceName),
		url.PathEscape(registryID),
	)
	res := &CloudProjectContainerRegistry{}

	log.Printf("[DEBUG] Will read registry %s and project: %s", registryID, serviceName)

	err := config.OVHClient.Get(endpoint, res)
	if err != nil {
		return fmt.Errorf("calling get %s %w", endpoint, err)
	}

	for k, v := range res.ToMap() {
		if k == "iam_enabled" {
			d.Set(k, v)
		}
	}

	d.SetId(serviceName + "/" + registryID)

	log.Printf("[DEBUG] Read registry %+v", res)
	return nil
}

func resourceCloudProjectContainerRegistryIAMCreate(d *schema.ResourceData, meta any) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	registryID := d.Get("registry_id").(string)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s/iam",
		url.PathEscape(serviceName),
		url.PathEscape(registryID),
	)
	params := (&CloudProjectContainerRegistryIAMCreateOpts{}).FromResource(d)

	log.Printf("[DEBUG] Will enable registry %s IAM: %+v", registryID, params)

	err := config.OVHClient.Post(endpoint, params, nil)
	if err != nil {
		return fmt.Errorf("calling Post %s with params %v:\n\t %w", endpoint, params, err)
	}

	d.SetId(serviceName + "/" + registryID)
	d.Set("iam_enabled", true)

	log.Printf("[DEBUG] Registry %s IAM enabled", registryID)

	return nil
}

func resourceCloudProjectContainerRegistryIAMDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	registryID := d.Get("registry_id").(string)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s/iam",
		url.PathEscape(serviceName),
		url.PathEscape(registryID),
	)

	log.Printf("[DEBUG] Will disable registry %s IAM", registryID)
	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return fmt.Errorf("calling delete %s %w", endpoint, err)
	}

	log.Printf("[DEBUG] Registry %s IAM disabled", registryID)

	return nil
}

package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudProjectContainerRegistryIAM() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectContainerRegistryIAMRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "Service name",
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"registry_id": {
				Type:        schema.TypeString,
				Description: "Registry ID",
				Required:    true,
				ForceNew:    true,
			},
			"iam_enabled": {
				Type:        schema.TypeBool,
				Description: "OVHcloud IAM enabled",
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudProjectContainerRegistryIAMRead(d *schema.ResourceData, meta any) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	registryID := d.Get("registry_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/containerRegistry/%s", serviceName, registryID)
	res := &CloudProjectContainerRegistry{}

	log.Printf("[DEBUG] Will read from registry %s and project: %s", registryID, serviceName)

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

	log.Printf("[DEBUG] Read IAM %+v", res)

	return nil
}

package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudProjectContainerRegistry() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectContainerRegistryRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"registry_id": {
				Type:        schema.TypeString,
				Description: "Registry ID",
				Required:    true,
			},

			// Computed
			"created_at": {
				Type:        schema.TypeString,
				Description: "Registry creation date",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Registry name",
				Computed:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "Project ID of your registry",
				Computed:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "Region of the registry.",
				Computed:    true,
			},
			"size": {
				Type:        schema.TypeInt,
				Description: "Current size of the registry (bytes)",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Registry status",
				Computed:    true,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Description: "Registry last update date",
				Computed:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "Access url of the registry",
				Computed:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "Version of your registry",
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudProjectContainerRegistryRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	id := d.Get("registry_id").(string)

	log.Printf("[DEBUG] Will read cloud project registry %s for project: %s", id, serviceName)

	reg := &CloudProjectContainerRegistry{}

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s",
		url.PathEscape(serviceName),
		url.PathEscape(id),
	)
	err := config.OVHClient.Get(endpoint, reg)
	if err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	for k, v := range reg.ToMap() {
		if k != "id" {
			d.Set(k, v)
		}
	}

	d.SetId(id)

	return nil
}

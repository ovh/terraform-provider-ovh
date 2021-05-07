package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudProjectCapabilitiesContainerRegistryFilter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectCapabilitiesContainerRegistryFilterRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"region": {
				Type:        schema.TypeString,
				Description: "Region of the registry.",
				Required:    true,
			},
			"plan_name": {
				Type:        schema.TypeString,
				Description: "Plan name of the registry.",
				Required:    true,
			},

			//Computed
			"code": {
				Type:        schema.TypeString,
				Description: "Plan code from catalog",
				Computed:    true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Description: "Plan creation date",
				Computed:    true,
			},
			"features": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Features of the plan",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vulnerability": {
							Type:        schema.TypeBool,
							Description: "Vulnerability scanning",
							Computed:    true,
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Plan name",
				Computed:    true,
			},
			"registry_limits": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Container registry limits",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_storage": {
							Type:        schema.TypeInt,
							Description: "Docker image storage limits in bytes",

							Computed: true,
						},
						"parallel_request": {
							Type:        schema.TypeInt,
							Description: "Parallel requests on Docker image API (/v2 Docker registry API)",
							Computed:    true,
						},
					},
				},
			},
			"updated_at": {
				Type:        schema.TypeString,
				Description: "Plan last update date",
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudProjectCapabilitiesContainerRegistryFilterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	log.Printf("[DEBUG] Will read cloud project capabilities container registry for project: %s", serviceName)

	capregs := []CloudProjectCapabilitiesContainerRegistry{}

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/capabilities/containerRegistry",
		url.PathEscape(serviceName),
	)
	err := config.OVHClient.Get(endpoint, &capregs)
	if err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	match := false
	for _, capreg := range capregs {
		if capreg.RegionName == d.Get("region").(string) {
			for _, plan := range capreg.Plans {
				if plan.Name == d.Get("plan_name").(string) {
					for k, v := range plan.ToMap() {
						match = true
						if k == "id" {
							d.SetId(v.(string))
						} else {
							d.Set(k, v)
						}
					}
				}
			}
		}
	}

	if !match {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	return nil
}

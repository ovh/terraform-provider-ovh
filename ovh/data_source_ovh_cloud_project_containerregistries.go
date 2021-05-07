package ovh

import (
	"fmt"
	"log"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceCloudProjectContainerRegistries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectContainerRegistriesRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"result": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_at": {
							Type:        schema.TypeString,
							Description: "Registry creation date",
							Computed:    true,
						},
						"id": {
							Type:        schema.TypeString,
							Description: "Registry ID",
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
				},
			},
		},
	}
}

func dataSourceCloudProjectContainerRegistriesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	log.Printf("[DEBUG] Will read cloud project registries for project: %s", serviceName)

	regs := []CloudProjectContainerRegistry{}

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry",
		url.PathEscape(serviceName),
	)
	err := config.OVHClient.Get(endpoint, &regs)
	if err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	mapregs := make([]map[string]interface{}, len(regs))
	ids := make([]string, len(regs))

	for i, reg := range regs {
		mapregs[i] = reg.ToMap()
		ids = append(ids, reg.Id)
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(ids)

	d.SetId(hashcode.Strings(ids))
	d.Set("result", mapregs)

	return nil
}

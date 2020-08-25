package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudRegionsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_PROJECT_ID", nil),
			},
			"has_services_up": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"names": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceCloudRegionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	projectId := d.Get("project_id").(string)

	log.Printf("[DEBUG] Will read public cloud regions for project: %s", projectId)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/region",
		url.PathEscape(projectId),
	)

	names := make([]string, 0)
	err := config.OVHClient.Get(endpoint, &names)

	if err != nil {
		return fmt.Errorf("Error calling %s:\n\t %q", endpoint, err)
	}

	d.SetId(projectId)

	var services []interface{}
	if servicesVal, ok := d.GetOk("has_services_up"); ok {
		services = servicesVal.(*schema.Set).List()
	}

	// no filtering on services
	if len(services) < 1 {
		d.Set("names", names)
		return nil
	}

	filtered_names := make([]string, 0)
	for _, n := range names {
		region, err := getCloudRegion(projectId, n, config.OVHClient)
		if err != nil {
			return err
		}

		for _, service := range services {
			if region.HasServiceUp(service.(string)) {
				filtered_names = append(filtered_names, n)
			}
		}
	}

	d.Set("names", filtered_names)

	return nil
}

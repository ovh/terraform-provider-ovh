package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudProjectRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectRegionsRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
				Description: "Service name of the resource representing the id of the cloud project.",
			},
			"has_services_up": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceCloudProjectRegionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	log.Printf("[DEBUG] Will read public cloud regions for project: %s", serviceName)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/region",
		url.PathEscape(serviceName),
	)

	names := make([]string, 0)
	if err := config.OVHClient.Get(endpoint, &names); err != nil {
		return fmt.Errorf("Error calling %s:\n\t %q", endpoint, err)
	}

	d.SetId(serviceName)
	d.Set("service_name", serviceName)

	var services []interface{}
	if servicesVal, ok := d.GetOk("has_services_up"); ok {
		services = servicesVal.([]interface{})
	}

	// no filtering on services
	if len(services) < 1 {
		d.Set("names", names)
		return nil
	}

	filtered_names := make([]string, 0)
	for _, n := range names {
		region, err := getCloudProjectRegion(serviceName, n, config.OVHClient)
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

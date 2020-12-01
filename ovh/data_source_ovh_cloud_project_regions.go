package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func dataSourceCloudProjectRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectRegionsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				DefaultFunc:   schema.EnvDefaultFunc("OVH_PROJECT_ID", nil),
				Description:   "Id of the cloud project. DEPRECATED, use `service_name` instead",
				ConflictsWith: []string{"service_name"},
			},
			"service_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				DefaultFunc:   schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
				Description:   "Service name of the resource representing the id of the cloud project.",
				ConflictsWith: []string{"project_id"},
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
	serviceName, err := helpers.GetCloudProjectServiceName(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Will read public cloud regions for project: %s", serviceName)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/region",
		url.PathEscape(serviceName),
	)

	names := make([]string, 0)
	err = config.OVHClient.Get(endpoint, &names)

	if err != nil {
		return fmt.Errorf("Error calling %s:\n\t %q", endpoint, err)
	}

	d.SetId(serviceName)
	d.Set("service_name", serviceName)
	d.Set("project_id", serviceName)

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

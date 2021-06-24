package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceCloudProjectRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectRegionRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
				Description: "Service name of the resource representing the id of the cloud project.",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"services": {
				Type:     schema.TypeSet,
				Set:      cloudServiceHash,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"continent_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datacenter_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudProjectRegionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	name := d.Get("name").(string)

	log.Printf("[DEBUG] Will read public cloud region %s for project: %s", name, serviceName)

	region, err := getCloudProjectRegion(serviceName, name, config.OVHClient)
	if err != nil {
		return err
	}

	d.Set("datacenter_location", region.DatacenterLocation)
	d.Set("continent_code", region.ContinentCode)

	services := &schema.Set{
		F: cloudServiceHash,
	}
	for i := range region.Services {
		service := map[string]interface{}{
			"name":   region.Services[i].Name,
			"status": region.Services[i].Status,
		}
		services.Add(service)
	}

	d.Set("services", services)
	d.Set("service_name", serviceName)
	d.SetId(fmt.Sprintf("%s_%s", serviceName, name))

	return nil
}

func getCloudProjectRegion(serviceName, region string, client *ovh.Client) (*CloudProjectRegionResponse, error) {
	log.Printf("[DEBUG] Will read public cloud region %s for project: %s", region, serviceName)

	response := &CloudProjectRegionResponse{}
	endpoint := fmt.Sprintf(
		"/cloud/project/%s/region/%s",
		url.PathEscape(serviceName),
		url.PathEscape(region),
	)
	err := client.Get(endpoint, response)

	if err != nil {
		return nil, fmt.Errorf("Error calling %s:\n\t %q", endpoint, err)
	}
	return response, nil
}

func cloudServiceHash(v interface{}) int {
	r := v.(map[string]interface{})
	return hashcode.String(r["name"].(string))
}

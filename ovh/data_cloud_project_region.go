package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the region",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Openstack region status",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region type (localzone, region, region-3-az)",
			},
			"services": {
				Type:        schema.TypeSet,
				Set:         cloudServiceHash,
				Computed:    true,
				Description: "Information about the different components available in the region",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Service name",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Service status",
							Computed:    true,
						},
						"endpoint": {
							Type:        schema.TypeString,
							Description: "Endpoint URL",
							Computed:    true,
						},
					},
				},
			},
			"continent_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region continent code",
			},
			"country_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region country code",
			},
			"datacenter_location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Location of the datacenter where the region is",
			},
			"availability_zones": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Availability zones of the region",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_countries": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Allowed countries for failover IP",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	d.Set("country_code", region.CountryCode)
	d.Set("status", region.Status)
	d.Set("type", region.Type)
	d.Set("ip_countries", region.IPCountries)
	d.Set("availability_zones", region.AvailabilityZones)

	services := &schema.Set{
		F: cloudServiceHash,
	}
	for _, service := range region.Services {
		services.Add(map[string]interface{}{
			"name":     service.Name,
			"status":   service.Status,
			"endpoint": service.Endpoint,
		})
	}

	d.Set("services", services)
	d.Set("service_name", serviceName)
	d.SetId(fmt.Sprintf("%s_%s", serviceName, name))

	return nil
}

func getCloudProjectRegion(serviceName, region string, client *OVHClient) (*CloudProjectRegionResponse, error) {
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

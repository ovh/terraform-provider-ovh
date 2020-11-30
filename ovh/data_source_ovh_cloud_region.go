package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceCloudRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudRegionRead,
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

			"continentCode": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Deprecated, use continent_code instead.",
			},
			"datacenterLocation": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Deprecated, use datacenter_location instead.",
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

func dataSourceCloudRegionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName, err := helpers.GetCloudProjectServiceName(d)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	log.Printf("[DEBUG] Will read public cloud region %s for project: %s", name, serviceName)

	region, err := getCloudRegion(serviceName, name, config.OVHClient)
	if err != nil {
		return err
	}

	// TODO: Deprecated - remove in next major release
	d.Set("datacenterLocation", region.DatacenterLocation)
	d.Set("continentCode", region.ContinentCode)
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
	d.Set("project_id", serviceName)
	d.SetId(fmt.Sprintf("%s_%s", serviceName, name))

	return nil
}

func getCloudRegion(serviceName, region string, client *ovh.Client) (*CloudRegionResponse, error) {
	log.Printf("[DEBUG] Will read public cloud region %s for project: %s", region, serviceName)

	response := &CloudRegionResponse{}
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

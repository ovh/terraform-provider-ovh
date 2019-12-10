package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/ovh/go-ovh/ovh"
)

func dataSourcePublicCloudRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePublicCloudRegionRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_PROJECT_ID", nil),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"services": {
				Type:     schema.TypeSet,
				Set:      publicCloudServiceHash,
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

func dataSourcePublicCloudRegionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	projectId := d.Get("project_id").(string)
	name := d.Get("name").(string)

	log.Printf("[DEBUG] Will read public cloud region %s for project: %s", name, projectId)

	region, err := getCloudRegion(projectId, name, config.OVHClient)
	if err != nil {
		return err
	}

	// TODO: Deprecated - remove in next major release
	d.Set("datacenterLocation", region.DatacenterLocation)
	d.Set("continentCode", region.ContinentCode)
	d.Set("datacenter_location", region.DatacenterLocation)
	d.Set("continent_code", region.ContinentCode)

	services := &schema.Set{
		F: publicCloudServiceHash,
	}
	for i := range region.Services {
		service := map[string]interface{}{
			"name":   region.Services[i].Name,
			"status": region.Services[i].Status,
		}
		services.Add(service)
	}

	d.Set("services", services)
	d.SetId(fmt.Sprintf("%s_%s", projectId, name))

	return nil
}

func getCloudRegion(projectId, region string, client *ovh.Client) (*PublicCloudRegionResponse, error) {
	log.Printf("[DEBUG] Will read public cloud region %s for project: %s", region, projectId)

	response := &PublicCloudRegionResponse{}
	endpoint := fmt.Sprintf(
		"/cloud/project/%s/region/%s",
		url.PathEscape(projectId),
		url.PathEscape(region),
	)
	err := client.Get(endpoint, response)

	if err != nil {
		return nil, fmt.Errorf("Error calling %s:\n\t %q", endpoint, err)
	}
	return response, nil
}

func publicCloudServiceHash(v interface{}) int {
	r := v.(map[string]interface{})
	return hashcode.String(r["name"].(string))
}

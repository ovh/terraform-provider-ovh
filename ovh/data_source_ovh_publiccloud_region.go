package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
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
	d.Partial(true)

	response := &PublicCloudRegionResponse{}
	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s", projectId, name)
	err := config.OVHClient.Get(endpoint, response)

	if err != nil {
		return fmt.Errorf("Error calling %s:\n\t %q", endpoint, err)
	}

	// TODO: Deprecated - remove in next major release
	d.Set("datacenterLocation", response.DatacenterLocation)
	d.Set("continentCode", response.ContinentCode)

	d.Set("datacenter_location", response.DatacenterLocation)
	d.Set("continent_code", response.ContinentCode)

	services := &schema.Set{
		F: publicCloudServiceHash,
	}
	for i := range response.Services {
		service := map[string]interface{}{
			"name":   response.Services[i].Name,
			"status": response.Services[i].Status,
		}
		services.Add(service)
	}

	d.Set("services", services)

	d.Partial(false)
	d.SetId(fmt.Sprintf("%s_%s", projectId, name))

	return nil
}

func publicCloudServiceHash(v interface{}) int {
	r := v.(map[string]interface{})
	return hashcode.String(r["name"].(string))
}

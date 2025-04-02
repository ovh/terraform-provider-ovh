package ovh

import (
	"fmt"
	"log"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

func dataSourceCloudProjectInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectInstancesRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"region": {
				Type:        schema.TypeString,
				Description: "Instance region",
				Required:    true,
			},
			// computed
			"instances": {
				Type:        schema.TypeSet,
				Description: "List of instances",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"addresses": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Instance IP addresses",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:        schema.TypeString,
										Description: "IP address",
										Computed:    true,
									},
									"version": {
										Type:        schema.TypeInt,
										Description: "IP version",
										Computed:    true,
									},
								},
							},
						},
						"attached_volumes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Volumes attached to the instance",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Description: "Volume id",
										Computed:    true,
									},
								},
							},
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Description: "Availability Zone",
							Computed:    true,
						},
						"flavor_id": {
							Type:        schema.TypeString,
							Description: "Flavor id",
							Computed:    true,
						},
						"flavor_name": {
							Type:        schema.TypeString,
							Description: "Flavor name",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Instance name",
							Computed:    true,
						},
						"id": {
							Type:        schema.TypeString,
							Description: "Instance id",
							Computed:    true,
						},
						"image_id": {
							Type:        schema.TypeString,
							Description: "Image id",
							Computed:    true,
						},
						"ssh_key": {
							Type:        schema.TypeString,
							Description: "SSH Key pair name",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Instance status",
							Computed:    true,
						},
						"task_state": {
							Type:        schema.TypeString,
							Description: "Instance task state",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudProjectInstancesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	region := d.Get("region").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/instance",
		url.PathEscape(serviceName),
		url.PathEscape(region),
	)
	var res []CloudProjectInstanceResponse

	log.Printf("[DEBUG] Will read instances from project %s in region %s", serviceName, region)
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	instances := make([]map[string]interface{}, 0, len(res))
	ids := make([]string, 0, len(res))
	for _, instance := range res {
		instances = append(instances, instance.ToMap())
		ids = append(ids, instance.Id)
	}
	sort.Strings(ids)

	d.SetId(hashcode.Strings(ids))
	d.Set("instances", instances)

	log.Printf("[DEBUG] Read instances: %s", ids)
	return nil
}

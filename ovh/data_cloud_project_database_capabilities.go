package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func dataSourceCloudProjectDatabaseCapabilities() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectDatabaseCapabilitiesRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},

			// Computed
			"engines": {
				Type:        schema.TypeSet,
				Description: "Database engines available",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_version": {
							Type:        schema.TypeString,
							Description: "Default version used for the engine",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Description of the engine",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Engine name",
							Computed:    true,
						},
						"ssl_modes": {
							Type:        schema.TypeSet,
							Description: "SSL modes for this engine",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"versions": {
							Type:        schema.TypeSet,
							Description: "Versions available for this engine",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"flavors": {
				Type:        schema.TypeSet,
				Description: "Flavors available",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"core": {
							Type:        schema.TypeInt,
							Description: "Flavor core number",
							Computed:    true,
						},
						"memory": {
							Type:        schema.TypeInt,
							Description: "Flavor ram size in GB",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the flavor",
							Computed:    true,
						},
						"storage": {
							Type:        schema.TypeInt,
							Description: "Flavor disk size in GB",
							Computed:    true,
						},
					},
				},
			},
			"options": {
				Type:        schema.TypeSet,
				Description: "Options available",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the option",
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Type of the option",
							Computed:    true,
						},
					},
				},
			},
			"plans": {
				Type:        schema.TypeSet,
				Description: "Plans available",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_retention": {
							Type:        schema.TypeString,
							Description: "Automatic backup retention duration",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Description of the plan",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the plan",
							Computed:    true,
						},
					},
				},
			},
			"availability": {
				Type:        schema.TypeSet,
				Description: "Availability of databases engines on cloud projects",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup": {
							Type:        schema.TypeString,
							Description: "Defines the type of backup",
							Computed:    true,
						},
						"default": {
							Type:        schema.TypeBool,
							Description: "Whether this availability can be used by default",
							Computed:    true,
						},
						"end_of_life": {
							Type:        schema.TypeString,
							Description: "End of life of the product",
							Computed:    true,
						},
						"engine": {
							Type:        schema.TypeString,
							Description: "Database engine name",
							Computed:    true,
						},
						"flavor": {
							Type:        schema.TypeString,
							Description: "Flavor name",
							Computed:    true,
						},
						"max_disk_size": {
							Type:        schema.TypeInt,
							Description: "Maximum possible disk size in GB",
							Computed:    true,
						},
						"max_node_number": {
							Type:        schema.TypeInt,
							Description: "Maximum nodes of the cluster",
							Computed:    true,
						},
						"min_disk_size": {
							Type:        schema.TypeInt,
							Description: "Minimum possible disk size in GB",
							Computed:    true,
						},
						"min_node_number": {
							Type:        schema.TypeInt,
							Description: "Minimum nodes of the cluster",
							Computed:    true,
						},
						"network": {
							Type:        schema.TypeString,
							Description: "Type of network",
							Computed:    true,
						},
						"plan": {
							Type:        schema.TypeString,
							Description: "Plan name",
							Computed:    true,
						},
						"region": {
							Type:        schema.TypeString,
							Description: "Region name",
							Computed:    true,
						},
						"start_date": {
							Type:        schema.TypeString,
							Description: "Date of the release of the product",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status of the availability",
							Computed:    true,
						},
						"step_disk_size": {
							Type:        schema.TypeInt,
							Description: "Flex disk size step in GB",
							Computed:    true,
						},
						"upstream_end_of_life": {
							Type:        schema.TypeString,
							Description: "End of life of the upstream product",
							Computed:    true,
						},
						"version": {
							Type:        schema.TypeString,
							Description: "Version name",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudProjectDatabaseCapabilitiesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	capabilitiesEndpoint := fmt.Sprintf("/cloud/project/%s/database/capabilities",
		url.PathEscape(serviceName),
	)
	capabilitiesRes := &CloudProjectDatabaseCapabilitiesResponse{}

	log.Printf("[DEBUG] Will read capabilities from project %s", serviceName)
	if err := config.OVHClient.Get(capabilitiesEndpoint, capabilitiesRes); err != nil {
		return helpers.CheckDeleted(d, err, capabilitiesEndpoint)
	}

	availabilityEndpoint := fmt.Sprintf("/cloud/project/%s/database/availability",
		url.PathEscape(serviceName),
	)
	availabilityRes := make([]CloudProjectDatabaseAvailability, 0)

	log.Printf("[DEBUG] Will read availability from project %s", serviceName)
	if err := config.OVHClient.Get(availabilityEndpoint, &availabilityRes); err != nil {
		return helpers.CheckDeleted(d, err, availabilityEndpoint)
	}

	capabilitiesRes.Availability = availabilityRes

	d.SetId(serviceName)
	for k, v := range capabilitiesRes.ToMap() {
		d.Set(k, v)
	}

	log.Printf("[DEBUG] Read capabilities %+v", capabilitiesRes)
	return nil
}

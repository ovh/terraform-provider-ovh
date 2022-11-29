package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func dataSourceCloudProjectDatabaseDatabase() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectDatabaseDatabaseRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"engine": {
				Type:         schema.TypeString,
				Description:  "Name of the engine of the service",
				Required:     true,
				ValidateFunc: helpers.ValidateEnum([]string{"mysql", "postgresql"}),
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Cluster ID",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the database",
				Required:    true,
			},

			//Computed
			"default": {
				Type:        schema.TypeBool,
				Description: "Defines if the database has been created by default",
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudProjectDatabaseDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	listEndpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/database",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterId),
	)

	listRes := make([]string, 0)

	log.Printf("[DEBUG] Will read databases from cluster %s from project %s", clusterId, serviceName)
	if err := config.OVHClient.Get(listEndpoint, &listRes); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", listEndpoint, err)
	}

	for _, id := range listRes {
		endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/database/%s",
			url.PathEscape(serviceName),
			url.PathEscape(engine),
			url.PathEscape(clusterId),
			url.PathEscape(id),
		)

		res := &CloudProjectDatabaseDatabaseResponse{}

		log.Printf("[DEBUG] Will read database %s from cluster %s from project %s", id, clusterId, serviceName)
		if err := config.OVHClient.Get(endpoint, res); err != nil {
			return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
		}

		if res.Name == name {
			for k, v := range res.ToMap() {
				if k != "id" {
					d.Set(k, v)
				} else {
					d.SetId(fmt.Sprint(v))
				}
			}
			log.Printf("[DEBUG] Read database %+v", res)
			return nil
		}
	}

	return fmt.Errorf("Database name %s not found for cluster %s from project %s", name, clusterId, serviceName)
}

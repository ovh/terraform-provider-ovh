package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudProjectDatabaseM3dbNamespace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectDatabaseM3dbNamespaceRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Cluster ID",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the namespace",
				Required:    true,
			},

			//Computed
			"resolution": {
				Type:        schema.TypeString,
				Description: "Resolution for an aggregated namespace",
				Computed:    true,
			},
			"retention_block_data_expiration_duration": {
				Type:        schema.TypeString,
				Description: "Controls how long we wait before expiring stale data",
				Computed:    true,
			},
			"retention_block_size_duration": {
				Type:        schema.TypeString,
				Description: "Controls how long to keep a block in memory before flushing to a fileset on disk",
				Computed:    true,
			},
			"retention_buffer_future_duration": {
				Type:        schema.TypeString,
				Description: "Controls how far into the future writes to the namespace will be accepted",
				Computed:    true,
			},
			"retention_buffer_past_duration": {
				Type:        schema.TypeString,
				Description: "Controls how far into the past writes to the namespace will be accepted",
				Computed:    true,
			},
			"retention_period_duration": {
				Type:        schema.TypeString,
				Description: "Controls the duration of time that M3DB will retain data for the namespace",
				Computed:    true,
			},
			"snapshot_enabled": {
				Type:        schema.TypeBool,
				Description: "Defines whether M3DB will create snapshot files for this namespace",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Type of namespace",
				Computed:    true,
			},
			"writes_to_commit_log_enabled": {
				Type:        schema.TypeBool,
				Description: "Defines whether M3DB will include writes to this namespace in the commit log",
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudProjectDatabaseM3dbNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)

	listEndpoint := fmt.Sprintf("/cloud/project/%s/database/m3db/%s/namespace",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
	)

	listRes := make([]string, 0)

	log.Printf("[DEBUG] Will read namespaces from cluster %s from project %s", clusterId, serviceName)
	if err := config.OVHClient.Get(listEndpoint, &listRes); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", listEndpoint, err)
	}

	name := d.Get("name").(string)
	for _, id := range listRes {
		endpoint := fmt.Sprintf("/cloud/project/%s/database/m3db/%s/namespace/%s",
			url.PathEscape(serviceName),
			url.PathEscape(clusterId),
			url.PathEscape(id),
		)
		res := &CloudProjectDatabaseM3dbNamespaceResponse{}

		log.Printf("[DEBUG] Will read namespace %s from cluster %s from project %s", id, clusterId, serviceName)
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
			log.Printf("[DEBUG] Read namespace %+v", res)
			return nil
		}
	}

	return fmt.Errorf("Namespace name %s not found for cluster %s from project %s", name, clusterId, serviceName)
}

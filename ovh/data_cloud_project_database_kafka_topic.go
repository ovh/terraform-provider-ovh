package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func dataSourceCloudProjectDatabaseKafkaTopic() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProjectDatabaseKafkaTopicRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Id of the database cluster",
				Required:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "Topic ID",
				Required:    true,
			},

			// Computed
			"min_insync_replicas": {
				Type:        schema.TypeInt,
				Description: "Minimum insync replica accepted for this topic",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the topic",
				Computed:    true,
			},
			"partitions": {
				Type:        schema.TypeInt,
				Description: "Number of partitions for this topic",
				Computed:    true,
			},
			"replication": {
				Type:        schema.TypeInt,
				Description: "Number of replication for this topic",
				Computed:    true,
			},
			"retention_bytes": {
				Type:        schema.TypeInt,
				Description: "Number of bytes for the retention of the data for this topic",
				Computed:    true,
			},
			"retention_hours": {
				Type:        schema.TypeInt,
				Description: "Number of hours for the retention of the data for this topic",
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudProjectDatabaseKafkaTopicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Get("id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/topic/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)
	res := &CloudProjectDatabaseKafkaTopicResponse{}

	log.Printf("[DEBUG] Will read topic %s from cluster %s from project %s", id, clusterId, serviceName)
	if err := config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	log.Printf("[DEBUG] Read topic %+v", res)
	return nil
}

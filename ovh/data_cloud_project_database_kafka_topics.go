package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceCloudProjectDatabaseKafkaTopics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProjectDatabaseKafkaTopicsRead,
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

			// Computed
			"topic_ids": {
				Type:        schema.TypeList,
				Description: "List of topic ids",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceCloudProjectDatabaseKafkaTopicsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/topic",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
	)
	res := make([]string, 0)

	log.Printf("[DEBUG] Will read topics from cluster %s from project %s", clusterId, serviceName)
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(res)

	d.SetId(hashcode.Strings(res))
	d.Set("topic_ids", res)
	return nil
}

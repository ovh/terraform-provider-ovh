package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func dataSourceCloudProjectDatabaseKafkaAcl() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectDatabaseKafkaAclRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Id of the database cluster",
				Required:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "Acl ID",
				Required:    true,
			},

			// Computed
			"permission": {
				Type:        schema.TypeString,
				Description: "Permission to give to this username on this topic",
				Computed:    true,
			},
			"topic": {
				Type:        schema.TypeString,
				Description: "Topic affected by this acl",
				Computed:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "Username affected by this acl",
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudProjectDatabaseKafkaAclRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Get("id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/acl/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)
	res := &CloudProjectDatabaseKafkaAclResponse{}

	log.Printf("[DEBUG] Will read acl %s from cluster %s from project %s", id, clusterId, serviceName)
	if err := config.OVHClient.Get(endpoint, res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	log.Printf("[DEBUG] Read acl %+v", res)
	return nil
}

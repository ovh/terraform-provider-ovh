package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func dataSourceCloudProjectDatabaseOpensearchPattern() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectDatabaseOpensearchPatternRead,
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
				Description: "Pattern ID",
				Required:    true,
			},

			// Computed
			"max_index_count": {
				Type:        schema.TypeInt,
				Description: "Maximum number of index for this pattern",
				Computed:    true,
			},
			"pattern": {
				Type:        schema.TypeString,
				Description: "Pattern format",
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudProjectDatabaseOpensearchPatternRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Get("id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/opensearch/%s/pattern/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)
	res := &CloudProjectDatabaseOpensearchPatternResponse{}

	log.Printf("[DEBUG] Will read pattern %s from cluster %s from project %s", id, clusterId, serviceName)
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

	log.Printf("[DEBUG] Read pattern %+v", res)
	return nil
}

package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceCloudProjectDatabaseKafkaUserAccess() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectDatabaseKafkaUserAccessRead,
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
			"user_id": {
				Type:        schema.TypeString,
				Description: "Id of the user",
				Required:    true,
			},

			// Computed
			"cert": {
				Type:        schema.TypeString,
				Description: "User cert",
				Computed:    true,
			},
			"key": {
				Type:        schema.TypeString,
				Description: "User key for the cert",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func dataSourceCloudProjectDatabaseKafkaUserAccessRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	userId := d.Get("user_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/user/%s/access",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(userId),
	)
	res := &CloudProjectDatabaseKafkaUserAccessResponse{}

	log.Printf("[DEBUG] Will read certificates of user %s from cluster %s from project %s", userId, clusterId, serviceName)
	if err := config.OVHClient.Get(endpoint, res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId(strconv.Itoa(hashcode.String(res.Cert)))
	for k, v := range res.ToMap() {
		d.Set(k, v)
	}

	log.Printf("[DEBUG] Read certificates %+v", res)
	return nil
}

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

func dataSourceCloudProjectDatabaseKafkaCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectDatabaseKafkaCertificatesRead,
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
			"ca": {
				Type:        schema.TypeString,
				Description: "CA certificate used for the service",
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudProjectDatabaseKafkaCertificatesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/certificates/",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
	)
	res := &CloudProjectDatabaseKafkaCertificatesResponse{}

	log.Printf("[DEBUG] Will read certificates from cluster %s from project %s", clusterId, serviceName)
	if err := config.OVHClient.Get(endpoint, res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId(strconv.Itoa(hashcode.String(res.Ca)))
	for k, v := range res.ToMap() {
		d.Set(k, v)
	}

	log.Printf("[DEBUG] Read certificates %+v", res)
	return nil
}

package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func dataSourceCloudProjectDatabaseMongodbPrometheus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProjectDatabaseMongodbPrometheusRead,
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

			//Computed
			"username": {
				Type:        schema.TypeString,
				Description: "Name of the user",
				Computed:    true,
			},
			"srv_domain": {
				Type:        schema.TypeString,
				Description: "Name of the srv domain endpoint",
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudProjectDatabaseMongodbPrometheusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterID := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/mongodb/%s/prometheus",
		url.PathEscape(serviceName),
		url.PathEscape(clusterID),
	)
	res := &CloudProjectDatabaseMongodbPrometheusEndpointResponse{}

	log.Printf("[DEBUG] Will read database %s from project: %s", d.Id(), serviceName)
	if err := config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	for k, v := range res.toMap() {
		d.Set(k, v)
	}
	d.SetId(clusterID)

	return nil
}

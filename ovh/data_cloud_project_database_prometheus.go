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

func dataSourceCloudProjectDatabasePrometheus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProjectDatabasePrometheusRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"engine": {
				Type:             schema.TypeString,
				Description:      "Name of the engine of the service",
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: helpers.ValidateDiagEnum([]string{"cassandra", "kafka", "kafkaConnect", "kafkaMirrorMaker", "mysql", "opensearch", "postgresql", "redis", "valkey"}),
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Id of the database cluster",
				ForceNew:    true,
				Required:    true,
			},

			//Computed
			"username": {
				Type:        schema.TypeString,
				Description: "Name of the user",
				Computed:    true,
			},
			"targets": {
				Type:        schema.TypeList,
				Description: "List of all endpoint targets",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Description: "Host of the endpoint",
							Computed:    true,
						},
						"port": {
							Type:        schema.TypeInt,
							Description: "Connection port for the endpoint",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudProjectDatabasePrometheusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterID := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/prometheus",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterID),
	)
	res := &CloudProjectDatabasePrometheusEndpointResponse{}

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

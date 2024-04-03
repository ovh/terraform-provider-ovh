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

func dataSourceCloudProjectDatabaseDatabase() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProjectDatabaseDatabaseRead,
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

func dataSourceCloudProjectDatabaseDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	listEndpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/database",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterID),
	)

	listRes := make([]string, 0)

	log.Printf("[DEBUG] Will read databases from cluster %s from project %s", clusterID, serviceName)
	if err := config.OVHClient.GetWithContext(ctx, listEndpoint, &listRes); err != nil {
		return diag.Errorf("Error calling GET %s:\n\t %q", listEndpoint, err)
	}

	for _, id := range listRes {
		endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/database/%s",
			url.PathEscape(serviceName),
			url.PathEscape(engine),
			url.PathEscape(clusterID),
			url.PathEscape(id),
		)

		res := &CloudProjectDatabaseDatabaseResponse{}

		log.Printf("[DEBUG] Will read database %s from cluster %s from project %s", id, clusterID, serviceName)
		if err := config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
			return diag.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
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

	return diag.Errorf("Database name %s not found for cluster %s from project %s", name, clusterID, serviceName)
}

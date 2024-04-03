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

func dataSourceCloudProjectDatabaseUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProjectDatabaseUserRead,
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
				ValidateFunc: helpers.ValidateEnum([]string{"cassandra", "mysql", "kafka", "kafkaConnect", "grafana"}),
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Cluster ID",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the user",
				Required:    true,
			},

			//Computed
			"created_at": {
				Type:        schema.TypeString,
				Description: "Date of the creation of the user",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Current status of the user",
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudProjectDatabaseUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	listEndpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/user",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterId),
	)

	listRes := make([]string, 0)

	log.Printf("[DEBUG] Will read users from cluster %s from project %s", clusterId, serviceName)
	if err := config.OVHClient.GetWithContext(ctx, listEndpoint, &listRes); err != nil {
		return diag.Errorf("Error calling GET %s:\n\t %q", listEndpoint, err)
	}

	for _, id := range listRes {
		endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/user/%s",
			url.PathEscape(serviceName),
			url.PathEscape(engine),
			url.PathEscape(clusterId),
			url.PathEscape(id),
		)

		res := &CloudProjectDatabaseUserResponse{}

		log.Printf("[DEBUG] Will read user %s from cluster %s from project %s", id, clusterId, serviceName)
		if err := config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
			return diag.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
		}

		if res.Username == name {
			for k, v := range res.ToMap() {
				if k != "id" {
					d.Set(k, v)
				} else {
					d.SetId(fmt.Sprint(v))
				}
			}
			log.Printf("[DEBUG] Read user %+v", res)
			return nil
		}
	}

	return diag.Errorf("User name %s not found for cluster %s from project %s", name, clusterId, serviceName)
}

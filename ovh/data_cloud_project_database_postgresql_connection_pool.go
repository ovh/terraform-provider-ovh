package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudProjectDatabasePostgresqlConnectionPool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProjectDatabasePostgresqlConnectionPoolRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Cluster ID",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the connection pool",
				Required:    true,
			},

			//Computed
			"database_id": {
				Type:        schema.TypeString,
				Description: "Database used for the connection pool",
				Computed:    true,
			},
			"mode": {
				Type:        schema.TypeString,
				Description: "Connection mode to the connection pool",
				Computed:    true,
			},
			"port": {
				Type:        schema.TypeInt,
				Description: "Port of the connection pool",
				Computed:    true,
			},
			"size": {
				Type:        schema.TypeInt,
				Description: "Size of the connection pool",
				Computed:    true,
			},
			"ssl_mode": {
				Type:        schema.TypeString,
				Description: "SSL connection mode for the pool",
				Computed:    true,
			},
			"uri": {
				Type:        schema.TypeString,
				Description: "Connection URI to the pool",
				Computed:    true,
			},
			"user_id": {
				Type:        schema.TypeString,
				Description: "User authorized to connect to the pool, if none all the users are allowed",
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudProjectDatabasePostgresqlConnectionPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	listEndpoint := fmt.Sprintf("/cloud/project/%s/database/postgresql/%s/connectionPool",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
	)

	listRes := make([]string, 0)

	log.Printf("[DEBUG] Will read connectionPools from cluster %s from project %s", clusterId, serviceName)
	if err := config.OVHClient.GetWithContext(ctx, listEndpoint, &listRes); err != nil {
		return diag.Errorf("Error calling GET %s:\n\t %q", listEndpoint, err)
	}

	for _, id := range listRes {
		endpoint := fmt.Sprintf("/cloud/project/%s/database/postgresql/%s/connectionPool/%s",
			url.PathEscape(serviceName),
			url.PathEscape(clusterId),
			url.PathEscape(id),
		)
		res := &CloudProjectDatabasePostgresqlConnectionPoolResponse{}

		log.Printf("[DEBUG] Will read connectionPool %s from cluster %s from project %s", id, clusterId, serviceName)
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
			log.Printf("[DEBUG] Read connectionPool %+v", res)
			return nil
		}
	}

	return diag.Errorf("ConnectionPool name %s not found for cluster %s from project %s", name, clusterId, serviceName)
}

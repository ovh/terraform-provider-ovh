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

func dataSourceCloudProjectDatabasePostgresqlConnectionPools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProjectDatabasePostgresqlConnectionPoolsRead,
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

			//Computed
			"connection_pool_ids": {
				Type:        schema.TypeList,
				Description: "List of connection pools ids",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceCloudProjectDatabasePostgresqlConnectionPoolsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterID := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/postgresql/%s/connectionPool",
		url.PathEscape(serviceName),
		url.PathEscape(clusterID),
	)
	res := make([]string, 0)

	log.Printf("[DEBUG] Will read connection pools from cluster %s from project %s", clusterID, serviceName)
	if err := config.OVHClient.GetWithContext(ctx, endpoint, &res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(res)

	d.SetId(hashcode.Strings(res))
	d.Set("connection_pool_ids", res)
	return nil
}

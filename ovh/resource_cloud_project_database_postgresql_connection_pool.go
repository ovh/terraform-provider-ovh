package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceCloudProjectDatabasePostgresqlConnectionPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudProjectDatabasePostgresqlConnectionPoolCreate,
		ReadContext:   resourceCloudProjectDatabasePostgresqlConnectionPoolRead,
		DeleteContext: resourceCloudProjectDatabasePostgresqlConnectionPoolDelete,
		UpdateContext: resourceCloudProjectDatabasePostgresqlConnectionPoolUpdate,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectDatabasePostgresqlConnectionPoolImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Id of the database cluster",
				ForceNew:    true,
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the connection pool",
				ForceNew:    true,
				Required:    true,
			},
			"database_id": {
				Type:        schema.TypeString,
				Description: "Database used for the connection pool",
				Required:    true,
			},
			"mode": {
				Type:        schema.TypeString,
				Description: "Connection mode to the connection pool",
				Required:    true,
			},
			"size": {
				Type:        schema.TypeInt,
				Description: "Size of the connection pool",
				Required:    true,
			},
			// Optional
			"user_id": {
				Type:        schema.TypeString,
				Description: "Database User authorized to connect to the pool, if none all the users are allowed",
				Optional:    true,
			},

			//Computed
			"port": {
				Type:        schema.TypeInt,
				Description: "Port of the connection pool",
				Computed:    true,
			},
			"ssl_mode": {
				Type:        schema.TypeString,
				Description: "Ssl connection mode for the pool",
				Computed:    true,
			},
			"uri": {
				Type:        schema.TypeString,
				Description: "Connection URI to the pool",
				Computed:    true,
			},
		},
	}
}

func resourceCloudProjectDatabasePostgresqlConnectionPoolImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return importCloudProjectDatabasePostgresqlConnectionPool(d, meta)
}

func resourceCloudProjectDatabasePostgresqlConnectionPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	f := func() interface{} {
		return (&CloudProjectDatabasePostgresqlConnectionPoolCreateOpts{}).FromResource(d)
	}
	return postCloudProjectDatabasePostgresqlConnectionPool(ctx, d, meta, dataSourceCloudProjectDatabasePostgresqlConnectionPoolRead, resourceCloudProjectDatabasePostgresqlConnectionPoolRead, resourceCloudProjectDatabasePostgresqlConnectionPoolUpdate, f)
}

func resourceCloudProjectDatabasePostgresqlConnectionPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/postgresql/%s/connectionPool/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)
	res := &CloudProjectDatabasePostgresqlConnectionPoolResponse{}

	log.Printf("[DEBUG] Will read connectionPool %s from cluster %s from project %s", id, clusterId, serviceName)
	if err := config.OVHClient.Get(endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

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

func resourceCloudProjectDatabasePostgresqlConnectionPoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	f := func() interface{} {
		return (&CloudProjectDatabasePostgresqlConnectionPoolUpdateOpts{}).FromResource(d)
	}
	return updateCloudProjectDatabasePostgresqlConnectionPool(ctx, d, meta, resourceCloudProjectDatabasePostgresqlConnectionPoolRead, f)
}

func resourceCloudProjectDatabasePostgresqlConnectionPoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return deleteCloudProjectDatabasePostgresqlConnectionPool(ctx, d, meta)
}

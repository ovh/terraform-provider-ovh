package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
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
				Type:         schema.TypeString,
				Description:  "Connection mode to the connection pool",
				Required:     true,
				ValidateFunc: helpers.ValidateEnum([]string{"session", "statement", "transaction"}),
			},
			"size": {
				Type:        schema.TypeInt,
				Description: "Size of the connection pool",
				Required:    true,
			},
			//Optional/Computed
			"user_id": {
				Type:        schema.TypeString,
				Description: "Database user authorized to connect to the pool, if none all the users are allowed",
				Optional:    true,
				Computed:    true,
			},

			//Computed
			"port": {
				Type:        schema.TypeInt,
				Description: "Port of the connection pool",
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
		},
	}
}

func resourceCloudProjectDatabasePostgresqlConnectionPoolImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenID := d.Id()
	n := 3
	splitID := strings.SplitN(givenID, "/", n)
	if len(splitID) != n {
		return nil, fmt.Errorf("import Id is not service_name/cluster_id/id formatted")
	}
	serviceName := splitID[0]
	clusterID := splitID[1]
	id := splitID[2]
	d.SetId(id)
	d.Set("cluster_id", clusterID)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectDatabasePostgresqlConnectionPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterID := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/postgresql/%s/connectionPool",
		url.PathEscape(serviceName),
		url.PathEscape(clusterID),
	)

	params := (&CloudProjectDatabasePostgresqlConnectionPoolCreateOpts{}).fromResource(d)
	res := &CloudProjectDatabasePostgresqlConnectionPoolResponse{}

	return diag.FromErr(
		retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
			func() *retry.RetryError {
				log.Printf("[DEBUG] Will create connection pool: %+v for cluster %s from project %s", params, clusterID, serviceName)
				rErr := config.OVHClient.PostWithContext(ctx, endpoint, params, res)
				if rErr != nil {
					// Manage a corner case where database is not create yet and connection pool POST return 403 Forbidden "Service database 'xxx' does not exist"
					if errOvh, ok := rErr.(*ovh.APIError); ok && (errOvh.Code == 403) {
						return retry.RetryableError(rErr)
					}
					return retry.NonRetryableError(rErr)
				}

				log.Printf("[DEBUG] Waiting for connection pool %s to be READY", res.ID)
				rErr = waitForCloudProjectDatabasePostgresqlConnectionPoolReady(ctx, config.OVHClient, serviceName, clusterID, res.ID, d.Timeout(schema.TimeoutCreate))
				if rErr != nil {
					return retry.NonRetryableError(fmt.Errorf("timeout while waiting connection pool %s to be READY: %s", res.ID, rErr.Error()))
				}
				log.Printf("[DEBUG] connection pool %s is READY", res.ID)

				d.SetId(res.ID)
				readDiags := resourceCloudProjectDatabasePostgresqlConnectionPoolRead(ctx, d, meta)
				rErr = diagnosticsToError(readDiags)
				if rErr != nil {
					return retry.NonRetryableError(rErr)
				}
				return nil
			},
		),
	)
}

func resourceCloudProjectDatabasePostgresqlConnectionPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterID := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/postgresql/%s/connectionPool/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterID),
		url.PathEscape(id),
	)
	res := &CloudProjectDatabasePostgresqlConnectionPoolResponse{}

	log.Printf("[DEBUG] Will read connectionPool %s from cluster %s from project %s", id, clusterID, serviceName)
	if err := config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	for k, v := range res.toMap() {
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
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterID := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/postgresql/%s/connectionPool/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterID),
		url.PathEscape(id),
	)

	params := (&CloudProjectDatabasePostgresqlConnectionPoolUpdateOpts{}).fromResource(d)

	log.Printf("[DEBUG] Will update connectionPool: %+v from cluster %s from project %s", params, clusterID, serviceName)
	err := config.OVHClient.PutWithContext(ctx, endpoint, params, nil)
	if err != nil {
		return diag.Errorf("calling Put %s with params %+v:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] connectionPool %s is READY", id)

	return resourceCloudProjectDatabasePostgresqlConnectionPoolRead(ctx, d, meta)
}

func resourceCloudProjectDatabasePostgresqlConnectionPoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterID := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/postgresql/%s/connectionPool/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterID),
		url.PathEscape(id),
	)

	log.Printf("[DEBUG] Will delete connectionPool %s from cluster %s from project %s", id, clusterID, serviceName)
	err := config.OVHClient.DeleteWithContext(ctx, endpoint, nil)
	if err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	log.Printf("[DEBUG] Waiting for connection pool %s to be DELETED", id)
	err = waitForCloudProjectDatabasePostgresqlConnectionPoolDeleted(ctx, config.OVHClient, serviceName, clusterID, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("timeout while waiting connection pool %s to be DELETED: %s", id, err.Error())
	}
	log.Printf("[DEBUG] connection pool %s is DELETED", id)

	d.SetId("")

	return nil
}

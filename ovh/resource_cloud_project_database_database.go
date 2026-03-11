package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceCloudProjectDatabaseDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudProjectDatabaseDatabaseCreate,
		ReadContext:   resourceCloudProjectDatabaseDatabaseRead,
		DeleteContext: resourceCloudProjectDatabaseDatabaseDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectDatabaseDatabaseImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

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
				ValidateDiagFunc: helpers.ValidateDiagEnum([]string{"clickhouse", "mysql", "postgresql"}),
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Id of the database cluster",
				ForceNew:    true,
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Database name",
				ForceNew:    true,
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

func resourceCloudProjectDatabaseDatabaseImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenID := d.Id()
	n := 4
	splitID := strings.SplitN(givenID, "/", n)
	if len(splitID) != n {
		return nil, fmt.Errorf("import Id is not service_name/engine/cluster_id/id formatted")
	}
	serviceName := splitID[0]
	engine := splitID[1]
	clusterID := splitID[2]
	id := splitID[3]
	d.SetId(id)
	d.Set("cluster_id", clusterID)
	d.Set("engine", engine)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectDatabaseDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterID := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/database",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterID),
	)

	params := (&CloudProjectDatabaseDatabaseCreateOpts{}).FromResource(d)
	res := &CloudProjectDatabaseDatabaseResponse{}

	log.Printf("[DEBUG] Will create database: %+v for cluster %s from project %s", params, clusterID, serviceName)
	err := config.OVHClient.PostWithContext(ctx, endpoint, params, res)
	if err != nil {
		return diag.Errorf("calling Post %s with params %+v:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for database %s to be READY", res.ID)
	err = waitForCloudProjectDatabaseDatabaseReady(ctx, config.OVHClient, serviceName, engine, clusterID, res.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("timeout while waiting database %s to be READY: %s", res.ID, err.Error())
	}
	log.Printf("[DEBUG] database %s is READY", res.ID)

	d.SetId(res.ID)

	return resourceCloudProjectDatabaseDatabaseRead(ctx, d, meta)
}

func resourceCloudProjectDatabaseDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterID := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/database/%s",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterID),
		url.PathEscape(id),
	)

	res := &CloudProjectDatabaseDatabaseResponse{}

	log.Printf("[DEBUG] Will read database %s from cluster %s from project %s", id, clusterID, serviceName)
	if err := config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	return nil
}

func resourceCloudProjectDatabaseDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterID := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/database/%s",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterID),
		url.PathEscape(id),
	)

	log.Printf("[DEBUG] Will delete database %s from cluster %s from project %s", id, clusterID, serviceName)
	err := config.OVHClient.DeleteWithContext(ctx, endpoint, nil)
	if err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	log.Printf("[DEBUG] Waiting for database %s to be DELETED", id)
	err = waitForCloudProjectDatabaseDatabaseDeleted(ctx, config.OVHClient, serviceName, engine, clusterID, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("timeout while waiting database %s to be DELETED: %s", id, err.Error())
	}
	log.Printf("[DEBUG] database %s is DELETED", id)

	d.SetId("")

	return nil
}

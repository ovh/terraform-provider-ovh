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
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceCloudProjectDatabaseM3dbNamespace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudProjectDatabaseM3dbNamespaceCreate,
		ReadContext:   resourceCloudProjectDatabaseM3dbNamespaceRead,
		DeleteContext: resourceCloudProjectDatabaseM3dbNamespaceDelete,
		UpdateContext: resourceCloudProjectDatabaseM3dbNamespaceUpdate,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectDatabaseM3dbNamespaceImportState,
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
				Description: "Name of the namespace",
				ForceNew:    true,
				Required:    true,
			},
			"resolution": {
				Type:             schema.TypeString,
				Description:      "Resolution for an aggregated namespace",
				Required:         true,
				DiffSuppressFunc: DiffDurationRfc3339,
			},
			"retention_block_data_expiration_duration": {
				Type:             schema.TypeString,
				Description:      "Controls how long we wait before expiring stale data",
				Optional:         true,
				DiffSuppressFunc: DiffDurationRfc3339,
			},
			"retention_buffer_future_duration": {
				Type:             schema.TypeString,
				Description:      "Controls how far into the future writes to the namespace will be accepted",
				Optional:         true,
				DiffSuppressFunc: DiffDurationRfc3339,
			},
			"retention_buffer_past_duration": {
				Type:             schema.TypeString,
				Description:      "Controls how far into the past writes to the namespace will be accepted",
				Optional:         true,
				DiffSuppressFunc: DiffDurationRfc3339,
			},
			"retention_period_duration": {
				Type:             schema.TypeString,
				Description:      "Controls the duration of time that M3DB will retain data for the namespace",
				Required:         true,
				DiffSuppressFunc: DiffDurationRfc3339,
			},
			"snapshot_enabled": {
				Type:        schema.TypeBool,
				Description: "Defines whether M3db will create snapshot files for this namespace",
				Optional:    true,
			},
			"writes_to_commit_log_enabled": {
				Type:        schema.TypeBool,
				Description: "Defines whether M3db will include writes to this namespace in the commit log",
				Optional:    true,
			},

			//Optional/Computed
			"retention_block_size_duration": {
				Type:             schema.TypeString,
				Description:      "Controls how long to keep a block in memory before flushing to a fileset on disk",
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: DiffDurationRfc3339,
			},

			// Computed
			"type": {
				Type:        schema.TypeString,
				Description: "Type of namespace",
				Computed:    true,
			},
		},
	}
}

func resourceCloudProjectDatabaseM3dbNamespaceImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	n := 3
	splitId := strings.SplitN(givenId, "/", n)
	if len(splitId) != n {
		return nil, fmt.Errorf("Import Id is not service_name/cluster_id/id formatted")
	}
	serviceName := splitId[0]
	clusterId := splitId[1]
	id := splitId[2]
	d.SetId(id)
	d.Set("cluster_id", clusterId)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectDatabaseM3dbNamespaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/m3db/%s/namespace",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
	)

	// Should read one time to
	listRes := make([]string, 0)
	log.Printf("[DEBUG] Will read namespaces from cluster %s from project %s", clusterId, serviceName)
	if err := config.OVHClient.Get(endpoint, &listRes); err != nil {
		return diag.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	params := (&CloudProjectDatabaseM3dbNamespaceCreateOpts{}).FromResource(d)
	res := &CloudProjectDatabaseM3dbNamespaceResponse{}

	log.Printf("[DEBUG] Will create namespace: %+v for cluster %s from project %s", params, clusterId, serviceName)
	err := config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		return diag.Errorf("calling Post %s with params %+v:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for namespace %s to be READY", res.Id)
	err = waitForCloudProjectDatabaseM3dbNamespaceReady(ctx, config.OVHClient, serviceName, clusterId, res.Id, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("timeout while waiting namespace %s to be READY: %s", res.Id, err.Error())
	}
	log.Printf("[DEBUG] namespace %s is READY", res.Id)

	d.SetId(res.Id)

	return resourceCloudProjectDatabaseM3dbNamespaceRead(ctx, d, meta)
}

func resourceCloudProjectDatabaseM3dbNamespaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/m3db/%s/namespace/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)
	res := &CloudProjectDatabaseM3dbNamespaceResponse{}

	log.Printf("[DEBUG] Will read namespace %s from cluster %s from project %s", id, clusterId, serviceName)
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

	return nil
}

func resourceCloudProjectDatabaseM3dbNamespaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/m3db/%s/namespace/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)
	params := (&CloudProjectDatabaseM3dbNamespaceUpdateOpts{}).FromResource(d)

	log.Printf("[DEBUG] Will update namespace: %+v from cluster %s from project %s", params, clusterId, serviceName)
	err := config.OVHClient.Put(endpoint, params, nil)
	if err != nil {
		return diag.Errorf("calling Put %s with params %+v:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for namespace %s to be READY", id)
	err = waitForCloudProjectDatabaseM3dbNamespaceReady(ctx, config.OVHClient, serviceName, clusterId, id, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.Errorf("timeout while waiting namespace %s to be READY: %s", id, err.Error())
	}
	log.Printf("[DEBUG] namespace %s is READY", id)

	return resourceCloudProjectDatabaseM3dbNamespaceRead(ctx, d, meta)
}

func resourceCloudProjectDatabaseM3dbNamespaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/m3db/%s/namespace/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)

	log.Printf("[DEBUG] Will delete namespace %s from cluster %s from project %s", id, clusterId, serviceName)
	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	log.Printf("[DEBUG] Waiting for namespace %s to be DELETED", id)
	err = waitForCloudProjectDatabaseM3dbNamespaceDeleted(ctx, config.OVHClient, serviceName, clusterId, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("timeout while waiting namespace %s to be DELETED: %s", id, err.Error())
	}
	log.Printf("[DEBUG] namespace %s is DELETED", id)

	d.SetId("")

	return nil
}

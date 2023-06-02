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

func resourceCloudProjectDatabaseIntegration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudProjectDatabaseIntegrationCreate,
		ReadContext:   resourceCloudProjectDatabaseIntegrationRead,
		DeleteContext: resourceCloudProjectDatabaseIntegrationDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectDatabaseIntegrationImportState,
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
				Type:        schema.TypeString,
				Description: "Name of the engine of the service",
				ForceNew:    true,
				Required:    true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value == "mongodb" {
						errors = append(errors, fmt.Errorf("Value %s is not a valid engine for integration", value))
					}
					return
				},
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Id of the database cluster",
				ForceNew:    true,
				Required:    true,
			},
			"destination_service_id": {
				Type:        schema.TypeString,
				Description: "ID of the destination service",
				ForceNew:    true,
				Required:    true,
			},
			"parameters": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Parameters for the integration",
				ForceNew:    true,
				Optional:    true,
			},
			"source_service_id": {
				Type:        schema.TypeString,
				Description: "ID of the source service",
				ForceNew:    true,
				Required:    true,
			},

			//Optional/Computed
			"type": {
				Type:         schema.TypeString,
				Description:  "Type of the integration",
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ValidateFunc: helpers.ValidateEnum([]string{"grafanaDashboard", "grafanaDatasource", "kafkaConnect", "kafkaLogs", "kafkaMirrorMaker", "m3aggregator", "m3dbMetrics", "opensearchLogs", "postgresqlMetrics"}),
			},

			//Computed
			"status": {
				Type:        schema.TypeString,
				Description: "Current status of the integration",
				Computed:    true,
			},
		},
	}
}

func resourceCloudProjectDatabaseIntegrationImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	n := 4
	splitId := strings.SplitN(givenId, "/", n)
	if len(splitId) != n {
		return nil, fmt.Errorf("Import Id is not service_name/engine/cluster_id/id formatted")
	}
	serviceName := splitId[0]
	engine := splitId[1]
	clusterId := splitId[2]
	id := splitId[3]
	d.SetId(id)
	d.Set("cluster_id", clusterId)
	d.Set("engine", engine)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectDatabaseIntegrationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterId := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/integration",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterId),
	)

	params := (&CloudProjectDatabaseIntegrationCreateOpts{}).FromResource(d)
	res := &CloudProjectDatabaseIntegrationResponse{}

	log.Printf("[DEBUG] Will create integration: %+v for cluster %s from project %s", params, clusterId, serviceName)
	err := config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		return diag.Errorf("calling Post %s with params %+v:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for integration %s to be READY", res.Id)
	err = waitForCloudProjectDatabaseIntegrationReady(ctx, config.OVHClient, serviceName, engine, clusterId, res.Id, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("timeout while waiting integration %s to be READY: %s", res.Id, err.Error())
	}
	log.Printf("[DEBUG] integration %s is READY", res.Id)

	d.SetId(res.Id)
	return resourceCloudProjectDatabaseIntegrationRead(ctx, d, meta)
}

func resourceCloudProjectDatabaseIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/integration/%s",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)

	res := &CloudProjectDatabaseIntegrationResponse{}

	log.Printf("[DEBUG] Will read integration %s from cluster %s from project %s", id, clusterId, serviceName)
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

func resourceCloudProjectDatabaseIntegrationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/integration/%s",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)

	log.Printf("[DEBUG] Will delete integration %s from cluster %s from project %s", id, clusterId, serviceName)
	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	log.Printf("[DEBUG] Waiting for integration %s to be DELETED", id)
	err = waitForCloudProjectDatabaseIntegrationDeleted(ctx, config.OVHClient, serviceName, engine, clusterId, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("timeout while waiting integration %s to be DELETED: %s", id, err.Error())
	}
	log.Printf("[DEBUG] integration %s is DELETED", id)

	d.SetId("")

	return nil
}

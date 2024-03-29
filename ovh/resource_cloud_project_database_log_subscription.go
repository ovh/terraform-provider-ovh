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

func resourceCloudProjectDatabaseLogSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudProjectDatabaseLogSubscriptionCreate,
		ReadContext:   resourceCloudProjectDatabaseLogSubscriptionRead,
		DeleteContext: resourceCloudProjectDatabaseLogSubscriptionDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectDatabaseLogSubscriptionImportState,
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
			"engine": {
				Type:        schema.TypeString,
				Description: "Name of the engine of the service",
				ForceNew:    true,
				Required:    true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Id of the database cluster",
				ForceNew:    true,
				Required:    true,
			},
			"stream_id": {
				Type:        schema.TypeString,
				Description: "Id of the target Log data platform stream",
				ForceNew:    true,
				Required:    true,
			},
			//Computed
			"created_at": {
				Type:        schema.TypeString,
				Description: "Creation date of the subscription",
				Computed:    true,
			},
			"kind": {
				Type:        schema.TypeString,
				Description: "Log kind name of this subscription",
				Computed:    true,
			},
			"ldp_service_name": {
				Type:        schema.TypeString,
				Description: "Name of the destination log service",
				Sensitive:   true,
				Computed:    true,
			},
			"operation_id": {
				Type:        schema.TypeString,
				Description: "Identifier of the operation",
				Computed:    true,
			},
			"resource_name": {
				Type:        schema.TypeString,
				Description: "Name of subscribed resource, where the logs come from",
				Computed:    true,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Description: "Type of subscribed resource, where the logs come from",
				Computed:    true,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Description: "Last update date of the subscription",
				Computed:    true,
			},
		},
	}
}

func resourceCloudProjectDatabaseLogSubscriptionImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

func resourceCloudProjectDatabaseLogSubscriptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterID := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/log/subscription",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterID),
	)
	params := (&CloudProjectDatabaseLogSubscriptionCreateOpts{}).fromResource(d)
	res := &CloudProjectDatabaseLogSubscriptionResponse{}

	log.Printf("[DEBUG] Will create Log subscrition : %+v for cluster %s from project %s", params, clusterID, serviceName)
	err := config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		diag.Errorf("calling Post %s with params %+v:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for Log subscription operation %s to be READY", res.OperationID)
	op, err := waitForDbaasLogsOperation(ctx, config.OVHClient, res.LDPServiceName, res.OperationID)
	if err != nil {
		return diag.Errorf("timeout while waiting log subscrition operation %s to be READY: %q", res.OperationID, err)
	}
	log.Printf("[DEBUG] Log subscrition operation %s is READY", res.OperationID)

	d.SetId(*op.SubscriptionID)

	readDiags := resourceCloudProjectDatabaseLogSubscriptionRead(ctx, d, meta)
	err = diagnosticsToError(readDiags)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceCloudProjectDatabaseLogSubscriptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterID := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/log/subscription/%s",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterID),
		url.PathEscape(id),
	)
	res := &CloudProjectDatabaseLogSubscriptionResponse{}

	log.Printf("[DEBUG] Will read log subscrition %s from cluster %s from project %s", id, clusterID, serviceName)
	if err := config.OVHClient.Get(endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	for k, v := range res.toMap() {
		if k == "operation_id" {
			continue
		} else if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	log.Printf("[DEBUG] Read log subscrition %+v", res)
	return nil
}

func resourceCloudProjectDatabaseLogSubscriptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterID := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/log/subscription/%s",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterID),
		url.PathEscape(id),
	)

	res := &CloudProjectDatabaseLogSubscriptionResponse{}

	log.Printf("[DEBUG] Will delete Log subscrition %s from cluster %s from project %s", id, clusterID, serviceName)
	err := config.OVHClient.Delete(endpoint, res)
	if err != nil {
		diag.Errorf("calling DELETE %s:\n\t %q", endpoint, err)
	}

	log.Printf("[DEBUG] Waiting for user %s to be DELETED", id)
	_, err = waitForDbaasLogsOperation(ctx, config.OVHClient, res.LDPServiceName, res.OperationID)
	if err != nil {
		return diag.Errorf("timeout while waiting log subscription %s to be DELETED: %q", id, err)
	}
	log.Printf("[DEBUG] Log subsription %s is DELETED", id)

	d.SetId("")

	return nil
}

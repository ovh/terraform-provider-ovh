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

func resourceCloudProjectRegionLoadbalancerLogSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudProjectRegionLoadbalancerSubscriptionsCreate,
		ReadContext:   resourceCloudProjectRegionLoadbalancerSubscriptionsRead,
		DeleteContext: resourceCloudProjectRegionLoadbalancerSubscriptionsDelete,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
				Description: "Service name of the resource representing the id of the cloud project.",
			},
			"region_name": {
				Type:        schema.TypeString,
				Description: "Region name of the resource representing the name of the region.",
				Required:    true,
				ForceNew:    true,
			},
			"loadbalancer_id": {
				Type:        schema.TypeString,
				Description: "ID representing the loadbalancer of the resource",
				Required:    true,
				ForceNew:    true,
			},
			"stream_id": {
				Type:        schema.TypeString,
				Description: "ID representing the stream of the resource",
				Required:    true,
				ForceNew:    true,
			},
			"kind": {
				Type:        schema.TypeString,
				Description: "Log kind name of this subscription",
				Required:    true,
				ForceNew:    true,
			},

			//computed
			"created_at": {
				Type:        schema.TypeString,
				Description: "Creation date of the subscription",
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
			"subscription_id": {
				Type:        schema.TypeString,
				Description: "Id of the subscription",
				Computed:    true,
			},
		},
	}
}

func resourceCloudProjectRegionLoadbalancerSubscriptionsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	regionName := d.Get("region_name").(string)
	loadbalancerID := d.Get("loadbalancer_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/loadbalancing/loadbalancer/%s/log/subscription",
		url.PathEscape(serviceName),
		url.PathEscape(regionName),
		url.PathEscape(loadbalancerID),
	)
	params := (&CloudProjectRegionLoadbalancerLogSubscriptionResourceCreateOpts{}).fromResource(d)
	res := &CreateCloudProjectRegionLoadbalancerLogSubscriptionResponse{}

	log.Printf("[DEBUG] Will create Log subscrition : %+v for loadbalancer %s on region %s from project %s", params, loadbalancerID, regionName, serviceName)
	err := config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		diag.Errorf("calling Post %s with params %+v:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for Log subscription operation %s to be READY", res.OperationID)
	op, err := waitForDbaasLogsOperation(ctx, config.OVHClient, res.ServiceName, res.OperationID)
	if err != nil {
		return diag.Errorf("timeout while waiting log subscrition operation %s to be READY: %q", res.ServiceName, err)
	}

	d.SetId(*op.SubscriptionID)

	return resourceCloudProjectRegionLoadbalancerSubscriptionsRead(ctx, d, meta)
}

func resourceCloudProjectRegionLoadbalancerSubscriptionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	regionName := d.Get("region_name").(string)
	loadbalancerID := d.Get("loadbalancer_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/loadbalancing/loadbalancer/%s/log/subscription/%s",
		url.PathEscape(serviceName),
		url.PathEscape(regionName),
		url.PathEscape(loadbalancerID),
		url.PathEscape(id),
	)
	res := &GetCloudProjectRegionLoadbalancerLogSubscriptionResponse{}

	log.Printf("[DEBUG] Will read log subscrition %s from loadbalancer %s on region %s from project %s", id, loadbalancerID, regionName, serviceName)
	if err := config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, fmt.Sprint(v))
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	log.Printf("[DEBUG] Read log subscrition %+v", res)
	return nil
}

func resourceCloudProjectRegionLoadbalancerSubscriptionsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	regionName := d.Get("region_name").(string)
	loadbalancerID := d.Get("loadbalancer_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/loadbalancing/loadbalancer/%s/log/subscription/%s",
		url.PathEscape(serviceName),
		url.PathEscape(regionName),
		url.PathEscape(loadbalancerID),
		url.PathEscape(id),
	)

	res := &GetCloudProjectRegionLoadbalancerLogSubscriptionResponse{}

	log.Printf("[DEBUG] Will delete Log subscrition for loadbalancer %s on region %s from project %s", loadbalancerID, regionName, serviceName)
	err := config.OVHClient.DeleteWithContext(ctx, endpoint, res)
	if err != nil {
		diag.Errorf("calling DELETE %s:\n\t %q", endpoint, err)
	}

	log.Printf("[DEBUG] Log subsription %s is DELETED", id)

	d.SetId("")

	return nil
}

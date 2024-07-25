package ovh

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"net/url"
)

func dataSourceCloudProjectRegionLoadbalancerLogSubscription() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProjectRegionLoadbalancerSubscriptionRead,
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
			"subscription_id": {
				Type:        schema.TypeString,
				Description: "ID representing the subscription",
				Required:    true,
				ForceNew:    true,
			},

			//computed
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
				//Sensitive:   true,
				Computed: true,
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
			"stream_id": {
				Type:        schema.TypeString,
				Description: "Id of the target Log data platform stream",
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

func dataSourceCloudProjectRegionLoadbalancerSubscriptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	regionName := d.Get("region_name").(string)
	loadbalancerID := d.Get("loadbalancer_id").(string)
	subscriptionID := d.Get("subscription_id").(string)

	log.Printf("[DEBUG] Will read public cloud loadbalancer %s log subscription %s for region %s for project: %s", loadbalancerID, subscriptionID, regionName, serviceName)

	response := &CloudProjectRegionLoadbalancerLogSubscriptionResponse{}
	endpoint := fmt.Sprintf(
		"/cloud/project/%s/region/%s/loadbalancing/loadbalancer/%s/log/subscription/%s",
		url.PathEscape(serviceName),
		url.PathEscape(regionName),
		url.PathEscape(loadbalancerID),
		url.PathEscape(subscriptionID),
	)

	if err := config.OVHClient.Get(endpoint, &response); err != nil {
		return diag.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	d.SetId(response.StreamID)
	d.Set("created_at", response.CreatedAt)
	d.Set("kind", response.Kind)
	d.Set("resource_name", response.Resource.Name)
	d.Set("resource_type", response.Resource.Type)
	d.Set("ldp_service_name", response.LDPServiceName)
	d.Set("stream_id", response.StreamID)
	d.Set("updated_at", response.UpdatedAt)
	return nil
}

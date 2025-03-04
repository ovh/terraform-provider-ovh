package ovh

import (
	"fmt"
	"log"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

func dataSourceCloudProjectRegionLoadbalancerLogSubscriptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectRegionLoadbalancerSubscriptionsRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
				Description: "Service name of the resource representing the id of the cloud project.",
			},
			"region_name": {
				Type:        schema.TypeString,
				Description: "Region name of the resource representing the name of the region.",
				Required:    true,
			},
			"loadbalancer_id": {
				Type:        schema.TypeString,
				Description: "ID representing the loadbalancer of the resource",
				Required:    true,
			},
			"kind": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Kind representing the loadbalancer.",
			},
			//computed
			"subscription_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceCloudProjectRegionLoadbalancerSubscriptionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	regionName := d.Get("region_name").(string)
	loadbalancerID := d.Get("loadbalancer_id").(string)
	kind := d.Get("kind").(string)
	query := ""
	if len(kind) > 0 {
		query = fmt.Sprintf("?kind=%s", url.PathEscape(kind))
	}
	log.Printf("[DEBUG] Will read public cloud loadbalancer %s log subscriptions for region %s for project: %s with query %s", loadbalancerID, regionName, serviceName, query)

	response := make([]string, 0)
	endpoint := fmt.Sprintf(
		"/cloud/project/%s/region/%s/loadbalancing/loadbalancer/%s/log/subscription%s",
		url.PathEscape(serviceName),
		url.PathEscape(regionName),
		url.PathEscape(loadbalancerID),
		query,
	)

	log.Printf("[DEBUG] Endpoint %s", query)

	if err := config.OVHClient.Get(endpoint, &response); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}
	sort.Strings(response)

	d.SetId(hashcode.Strings(response))
	d.Set("subscription_ids", response)
	return nil
}

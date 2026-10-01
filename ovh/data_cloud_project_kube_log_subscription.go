package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func dataSourceCloudProjectKubeLogSubscription() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProjectKubeLogSubscriptionRead,

		Schema: map[string]*schema.Schema{
			kubeServiceNameKey: {
				Type:        schema.TypeString,
				Description: "Service name of the resource representing the id of the cloud project.",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			kubeKubeIdKey: {
				Type:        schema.TypeString,
				Description: "Id of the managed kubernetes cluster.",
				Required:    true,
			},
			kubeLogSubscriptionIdKey: {
				Type:        schema.TypeString,
				Description: "Id of the subscription.",
				Required:    true,
			},
			// Computed
			kubeLogSubscriptionKindKey: {
				Type:        schema.TypeString,
				Description: "Log kind name of this subscription.",
				Computed:    true,
			},
			kubeLogSubscriptionStreamIdKey: {
				Type:        schema.TypeString,
				Description: "Id of the target Log data platform stream.",
				Computed:    true,
			},
			kubeCreatedAtKey: {
				Type:        schema.TypeString,
				Description: "Creation date of the subscription.",
				Computed:    true,
			},
			kubeUpdatedAtKey: {
				Type:        schema.TypeString,
				Description: "Last update date of the subscription.",
				Computed:    true,
			},
			kubeLogSubscriptionResourceKey: {
				Type:        schema.TypeList,
				Description: "Resource information of the subscription.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						kubeLogSubscriptionResourceNameKey: {
							Type:        schema.TypeString,
							Description: "Name of the subscribed resource.",
							Computed:    true,
						},
						kubeLogSubscriptionResourceTypeKey: {
							Type:        schema.TypeString,
							Description: "Type of the subscribed resource.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudProjectKubeLogSubscriptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get(kubeServiceNameKey).(string)
	kubeID := d.Get(kubeKubeIdKey).(string)
	subscriptionID := d.Get(kubeLogSubscriptionIdKey).(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/log/subscription/%s",
		url.PathEscape(serviceName),
		url.PathEscape(kubeID),
		url.PathEscape(subscriptionID),
	)
	res := &CloudProjectKubeLogSubscriptionResponse{}

	log.Printf("[DEBUG] Will read kube log subscription %s from kube %s from project %s", subscriptionID, kubeID, serviceName)
	if err := config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	for k, v := range res.toMap() {
		d.Set(k, v)
	}
	d.SetId(subscriptionID)

	log.Printf("[DEBUG] Read kube log subscription %+v", res)
	return nil
}

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

func resourceCloudProjectKubeLogSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudProjectKubeLogSubscriptionCreate,
		ReadContext:   resourceCloudProjectKubeLogSubscriptionRead,
		DeleteContext: resourceCloudProjectKubeLogSubscriptionDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectKubeLogSubscriptionImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			kubeServiceNameKey: {
				Type:        schema.TypeString,
				Description: "Service name of the resource representing the id of the cloud project.",
				ForceNew:    true,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			kubeKubeIdKey: {
				Type:        schema.TypeString,
				Description: "Id of the managed kubernetes cluster.",
				ForceNew:    true,
				Required:    true,
			},
			kubeLogSubscriptionKindKey: {
				Type:        schema.TypeString,
				Description: "Log kind name of this subscription.",
				ForceNew:    true,
				Required:    true,
			},
			kubeLogSubscriptionStreamIdKey: {
				Type:        schema.TypeString,
				Description: "Id of the target Log data platform stream.",
				ForceNew:    true,
				Required:    true,
			},
			// Computed
			kubeLogSubscriptionIdKey: {
				Type:        schema.TypeString,
				Description: "Id of the subscription.",
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

func resourceCloudProjectKubeLogSubscriptionImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenID := d.Id()
	n := 3
	splitID := strings.SplitN(givenID, "/", n)
	if len(splitID) != n {
		return nil, fmt.Errorf("import Id is not service_name/kube_id/subscription_id formatted")
	}
	serviceName := splitID[0]
	kubeID := splitID[1]
	subscriptionID := splitID[2]

	d.SetId(subscriptionID)
	d.Set(kubeServiceNameKey, serviceName)
	d.Set(kubeKubeIdKey, kubeID)
	d.Set(kubeLogSubscriptionIdKey, subscriptionID)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectKubeLogSubscriptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get(kubeServiceNameKey).(string)
	kubeID := d.Get(kubeKubeIdKey).(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/log/subscription",
		url.PathEscape(serviceName),
		url.PathEscape(kubeID),
	)
	params := (&CloudProjectKubeLogSubscriptionCreateOpts{}).fromResource(d)
	res := &CloudProjectKubeLogSubscriptionCreateResponse{}

	log.Printf("[DEBUG] Will create kube log subscription: %+v for kube %s from project %s", params, kubeID, serviceName)
	err := config.OVHClient.PostWithContext(ctx, endpoint, params, res)
	if err != nil {
		return diag.Errorf("calling Post %s with params %+v:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for kube log subscription operation %s to be READY", res.OperationId)
	op, err := waitForDbaasLogsOperation(ctx, config.OVHClient, res.ServiceName, res.OperationId)
	if err != nil {
		return diag.Errorf("timeout while waiting kube log subscription operation %s to be READY: %q", res.OperationId, err)
	}
	log.Printf("[DEBUG] Kube log subscription operation %s is READY", res.OperationId)

	d.SetId(*op.SubscriptionID)
	d.Set(kubeLogSubscriptionIdKey, *op.SubscriptionID)

	log.Printf("[DEBUG] Waiting for kube %s to be READY", kubeID)
	err = waitForCloudProjectKubeReady(config.OVHClient, serviceName, kubeID, []string{"REDEPLOYING", "UPDATING"}, []string{"READY"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("timeout while waiting for kube %s to be READY: %q", kubeID, err)
	}
	log.Printf("[DEBUG] kube %s is READY", kubeID)

	return resourceCloudProjectKubeLogSubscriptionRead(ctx, d, meta)
}

func resourceCloudProjectKubeLogSubscriptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get(kubeServiceNameKey).(string)
	kubeID := d.Get(kubeKubeIdKey).(string)
	subscriptionID := d.Id()

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

func resourceCloudProjectKubeLogSubscriptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get(kubeServiceNameKey).(string)
	kubeID := d.Get(kubeKubeIdKey).(string)
	subscriptionID := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/log/subscription/%s",
		url.PathEscape(serviceName),
		url.PathEscape(kubeID),
		url.PathEscape(subscriptionID),
	)

	log.Printf("[DEBUG] Will delete kube log subscription %s from kube %s from project %s", subscriptionID, kubeID, serviceName)
	err := config.OVHClient.DeleteWithContext(ctx, endpoint, nil)
	if err != nil {
		return diag.Errorf("calling DELETE %s:\n\t %q", endpoint, err)
	}

	log.Printf("[DEBUG] Deleted kube log subscription %s", subscriptionID)

	log.Printf("[DEBUG] Waiting for kube %s to be READY", kubeID)
	err = waitForCloudProjectKubeReady(config.OVHClient, serviceName, kubeID, []string{"REDEPLOYING", "UPDATING"}, []string{"READY"}, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("timeout while waiting for kube %s to be READY: %q", kubeID, err)
	}
	log.Printf("[DEBUG] kube %s is READY", kubeID)

	d.SetId("")

	return nil
}

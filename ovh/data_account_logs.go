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

func dataSourceAccountLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountLogsRead,

		Schema: map[string]*schema.Schema{
			"log_type": {
				Type:        schema.TypeString,
				Description: "Type of account logs to query (audit, activity, or access_policy)",
				Required:    true,
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Description: "ID of the subscription to retrieve",
				Required:    true,
			},
			// Computed attributes
			"created_at": {
				Type:        schema.TypeString,
				Description: "Creation date of the subscription",
				Computed:    true,
			},
			"kind": {
				Type:        schema.TypeString,
				Description: "Kind of log subscription",
				Computed:    true,
			},
			"ldp_service_name": {
				Type:        schema.TypeString,
				Description: "Name of the destination log service",
				Computed:    true,
				Sensitive:   true,
			},
			"stream_id": {
				Type:        schema.TypeString,
				Description: "ID of the target stream",
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

func dataSourceAccountLogsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	logType := d.Get("log_type").(string)
	subscriptionID := d.Get("subscription_id").(string)

	endpoint, err := getAccountLogsEndpoint(logType)
	if err != nil {
		return diag.FromErr(err)
	}

	endpoint = fmt.Sprintf("%s/%s", endpoint, url.PathEscape(subscriptionID))
	res := &AccountLogsResponse{}

	log.Printf("[DEBUG] Will read account logs subscription %s for log_type %s", subscriptionID, logType)
	if err := config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	d.SetId(subscriptionID)

	for k, v := range res.toMap() {
		if k != "id" {
			d.Set(k, v)
		}
	}

	log.Printf("[DEBUG] Read account logs subscription: %+v", res)
	return nil
}

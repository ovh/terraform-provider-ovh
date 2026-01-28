package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceAccountLogs() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccountLogsCreate,
		ReadContext:   resourceAccountLogsRead,
		UpdateContext: resourceAccountLogsUpdate,
		DeleteContext: resourceAccountLogsDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAccountLogsImportState,
		},

		Schema: map[string]*schema.Schema{
			"log_type": {
				Type:        schema.TypeString,
				Description: "Type of account logs to subscribe to (audit, activity, or access_policy)",
				ForceNew:    true,
				Required:    true,
			},
			"stream_id": {
				Type:        schema.TypeString,
				Description: "ID of the target Log Data Platform stream",
				ForceNew:    true,
				Required:    true,
			},
			"kind": {
				Type:        schema.TypeString,
				Description: "Kind of log subscription (default or other values)",
				ForceNew:    true,
				Required:    true,
			},
			// Computed attributes
			"created_at": {
				Type:        schema.TypeString,
				Description: "Creation date of the subscription",
				Computed:    true,
			},
			"ldp_service_name": {
				Type:        schema.TypeString,
				Description: "Name of the destination log service",
				Computed:    true,
				Sensitive:   true,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Description: "Last update date of the subscription",
				Computed:    true,
			},
		},
	}
}

func resourceAccountLogsImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenID := d.Id()
	// Format: logType/subscriptionId
	parts := strings.SplitN(givenID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("import ID should be formatted as logType/subscriptionId")
	}
	logType := parts[0]
	subscriptionID := parts[1]

	d.SetId(subscriptionID)
	d.Set("log_type", logType)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func getAccountLogsEndpoint(logType string) (string, error) {
	switch logType {
	case "audit":
		return "/me/logs/audit/log/subscription", nil
	case "activity":
		return "/me/api/log/subscription", nil
	case "access_policy":
		return "/iam/log/subscription", nil
	default:
		return "", fmt.Errorf("unsupported log type: %s", logType)
	}
}

func resourceAccountLogsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	logType := d.Get("log_type").(string)

	endpoint, err := getAccountLogsEndpoint(logType)
	if err != nil {
		return diag.FromErr(err)
	}

	params := (&AccountLogsCreateOpts{}).fromResource(d)
	res := &AccountLogsResponse{}

	log.Printf("[DEBUG] Will create account logs subscription for log_type %s", logType)
	if err := config.OVHClient.PostWithContext(ctx, endpoint, params, res); err != nil {
		return diag.Errorf("calling POST %s with params %+v:\n\t %q", endpoint, params, err)
	}

	d.SetId(res.SubscriptionID)
	d.Set("log_type", logType)

	log.Printf("[DEBUG] Account logs subscription created with ID: %s", res.SubscriptionID)

	return resourceAccountLogsRead(ctx, d, meta)
}

func resourceAccountLogsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	logType := d.Get("log_type").(string)
	subscriptionID := d.Id()

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

	for k, v := range res.toMap() {
		if k != "id" {
			d.Set(k, v)
		}
	}

	log.Printf("[DEBUG] Read account logs subscription: %+v", res)
	return nil
}

func resourceAccountLogsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Since all fields have ForceNew: true, updates are not supported
	// This is a placeholder for potential future enhancements
	log.Printf("[DEBUG] No updates available for account logs subscriptions")
	return nil
}

func resourceAccountLogsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	logType := d.Get("log_type").(string)
	subscriptionID := d.Id()

	endpoint, err := getAccountLogsEndpoint(logType)
	if err != nil {
		return diag.FromErr(err)
	}

	endpoint = fmt.Sprintf("%s/%s", endpoint, url.PathEscape(subscriptionID))

	log.Printf("[DEBUG] Will delete account logs subscription %s for log_type %s", subscriptionID, logType)
	if err := config.OVHClient.DeleteWithContext(ctx, endpoint, nil); err != nil {
		return diag.Errorf("calling DELETE %s:\n\t %q", endpoint, err)
	}

	log.Printf("[DEBUG] Account logs subscription %s deleted", subscriptionID)
	d.SetId("")

	return nil
}

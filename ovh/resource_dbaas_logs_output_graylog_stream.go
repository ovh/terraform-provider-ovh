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

func resourceDbaasLogsOutputGraylogStream() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDbaasLogsOutputGraylogStreamCreate,
		ReadContext:   resourceDbaasLogsOutputGraylogStreamRead,
		UpdateContext: resourceDbaasLogsOutputGraylogStreamUpdate,
		DeleteContext: resourceDbaasLogsOutputGraylogStreamDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDbaasLogsOutputGraylogStreamImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The service name",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Stream description",
				Required:    true,
			},
			"title": {
				Type:        schema.TypeString,
				Description: "Stream description",
				Required:    true,
			},

			// Optional ForceNew
			"parent_stream_id": {
				Type:        schema.TypeString,
				Description: "Parent stream ID",
				Optional:    true,
				ForceNew:    true,
			},
			"retention_id": {
				Type:        schema.TypeString,
				Description: "Retention ID",
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},

			// Optional
			"cold_storage_compression": {
				Type:        schema.TypeString,
				Description: "Cold storage compression method",
				Optional:    true,
				Computed:    true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateStringEnum(strings.ToUpper(v.(string)), []string{
						"LZMA",
						"GZIP",
						"DEFLATED",
						"ZSTD",
					})
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"cold_storage_content": {
				Type:        schema.TypeString,
				Description: "ColdStorage content",
				Optional:    true,
				Computed:    true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateStringEnum(strings.ToUpper(v.(string)), []string{
						"ALL",
						"GELF",
						"PLAIN",
					})
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"cold_storage_enabled": {
				Type:        schema.TypeBool,
				Description: "Is Cold storage enabled?",
				Computed:    true,
				Optional:    true,
			},
			"cold_storage_notify_enabled": {
				Type:        schema.TypeBool,
				Description: "Notify on new Cold storage archive",
				Computed:    true,
				Optional:    true,
			},
			"cold_storage_retention": {
				Type:        schema.TypeInt,
				Description: "Cold storage retention in year",
				Computed:    true,
				Optional:    true,
			},
			"cold_storage_target": {
				Type:        schema.TypeString,
				Description: "ColdStorage destination",
				Computed:    true,
				Optional:    true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateStringEnum(strings.ToUpper(v.(string)), []string{
						"PCA",
						"PCS",
					})
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"indexing_enabled": {
				Type:        schema.TypeBool,
				Description: "Enable ES indexing",
				Computed:    true,
				Optional:    true,
			},
			"indexing_max_size": {
				Type:        schema.TypeInt,
				Description: "Maximum indexing size (in GB)",
				Computed:    true,
				Optional:    true,
			},
			"indexing_notify_enabled": {
				Type:        schema.TypeBool,
				Description: "If set, notify when size is near 80, 90 or 100 % of the maximum configured setting",
				Computed:    true,
				Optional:    true,
			},
			"pause_indexing_on_max_size": {
				Type:        schema.TypeBool,
				Description: "If set, pause indexing when maximum size is reach",
				Computed:    true,
				Optional:    true,
			},
			"web_socket_enabled": {
				Type:        schema.TypeBool,
				Description: "Enable Websocket",
				Computed:    true,
				Optional:    true,
			},

			// computed
			"can_alert": {
				Type:        schema.TypeBool,
				Description: "Indicates if the current user can create alert on the stream",
				Computed:    true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Description: "Stream creation",
				Computed:    true,
			},
			"is_editable": {
				Type:        schema.TypeBool,
				Description: "Indicates if you are allowed to edit entry",
				Computed:    true,
			},
			"is_shareable": {
				Type:        schema.TypeBool,
				Description: "Indicates if you are allowed to share entry",
				Computed:    true,
			},
			"nb_alert_condition": {
				Type:        schema.TypeInt,
				Description: "Number of alert condition",
				Computed:    true,
			},
			"nb_archive": {
				Type:        schema.TypeInt,
				Description: "Number of coldstored archives",
				Computed:    true,
			},
			"stream_id": {
				Type:        schema.TypeString,
				Description: "Stream ID",
				Computed:    true,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Description: "Stream last update",
				Computed:    true,
			},
			"write_token": {
				Type:        schema.TypeString,
				Description: "Write token of the stream",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceDbaasLogsOutputGraylogStreamImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenID := d.Id()
	splitID := strings.SplitN(givenID, "/", 2)
	if len(splitID) != 2 {
		return nil, fmt.Errorf("Import Id is not service_name/id formatted")
	}
	serviceName := splitID[0]
	id := splitID[1]
	d.SetId(id)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceDbaasLogsOutputGraylogStreamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)

	log.Printf("[DEBUG] Will create dbaas logs output graylog stream for: %s", serviceName)

	opts := (&DbaasLogsOutputGraylogStreamCreateOpts{}).FromResource(d)
	res := &DbaasLogsOperation{}
	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/output/graylog/stream",
		url.PathEscape(serviceName),
	)
	if err := config.OVHClient.Post(endpoint, opts, res); err != nil {
		return diag.Errorf("Error calling post %s:\n\t %q", endpoint, err)
	}

	// Wait for operation status
	op, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, res.OperationId)
	if err != nil {
		return diag.FromErr(err)
	}

	id := op.StreamId
	if id == nil {
		return diag.Errorf("Stream Id is nil. This should not happen: operation is %s/%s", serviceName, res.OperationId)
	}

	d.SetId(*id)

	return resourceDbaasLogsOutputGraylogStreamRead(ctx, d, meta)
}

func resourceDbaasLogsOutputGraylogStreamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will update dbaas logs output graylog stream for: %s", serviceName)

	opts := (&DbaasLogsOutputGraylogStreamUpdateOpts{}).FromResource(d)
	res := &DbaasLogsOperation{}
	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/output/graylog/stream/%s",
		url.PathEscape(serviceName),
		url.PathEscape(id),
	)
	if err := config.OVHClient.Put(endpoint, opts, res); err != nil {
		return diag.Errorf("Error calling Put %s:\n\t %q", endpoint, err)
	}

	// Wait for operation status
	if _, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, res.OperationId); err != nil {
		return diag.FromErr(err)
	}

	return resourceDbaasLogsOutputGraylogStreamRead(ctx, d, meta)
}

func resourceDbaasLogsOutputGraylogStreamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will read dbaas logs output graylog stream: %s/%s", serviceName, id)
	res := &DbaasLogsOutputGraylogStream{}
	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/output/graylog/stream/%s",
		url.PathEscape(serviceName),
		url.PathEscape(id),
	)
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		log.Printf("[ERROR] %s: %v", endpoint, err)
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	// Get stream write token, if available
	writeToken, err := resourceDbaasLogsOutputGraylogStreamGetWriteToken(ctx, config, serviceName, id)
	if err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}
	d.Set("write_token", writeToken)

	return nil
}

func resourceDbaasLogsOutputGraylogStreamGetWriteToken(ctx context.Context, config *Config, serviceName, streamId string) (string, error) {
	var (
		ruleIds  []string
		endpoint = fmt.Sprintf("/dbaas/logs/%s/output/graylog/stream/%s/rule", url.PathEscape(serviceName), url.PathEscape(streamId))
	)

	if err := config.OVHClient.GetWithContext(ctx, endpoint, &ruleIds); err != nil {
		return "", fmt.Errorf("failed to list stream rules: %w", err)
	}

	for _, ruleId := range ruleIds {
		rule := DbaasLogsOutputGraylogStreamRule{}
		ruleEndpoint := endpoint + "/" + url.PathEscape(ruleId)

		if err := config.OVHClient.GetWithContext(ctx, ruleEndpoint, &rule); err != nil {
			return "", fmt.Errorf("failed to get stream rule %q: %w", ruleId, err)
		}

		if rule.Field == "X-OVH-TOKEN" {
			return rule.Value, nil
		}
	}

	return "", nil
}

func resourceDbaasLogsOutputGraylogStreamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will delete dbaas logs output graylog stream: %s/%s", serviceName, id)
	res := &DbaasLogsOperation{}
	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/output/graylog/stream/%s",
		url.PathEscape(serviceName),
		url.PathEscape(id),
	)

	if err := config.OVHClient.Delete(endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	// Wait for operation status
	if _, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, res.OperationId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

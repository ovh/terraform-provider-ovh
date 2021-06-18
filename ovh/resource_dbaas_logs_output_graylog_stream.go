package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceDbaasLogsOutputGraylogStream() *schema.Resource {
	return &schema.Resource{
		Create: resourceDbaasLogsOutputGraylogStreamCreate,
		Read:   resourceDbaasLogsOutputGraylogStreamRead,
		Update: resourceDbaasLogsOutputGraylogStreamUpdate,
		Delete: resourceDbaasLogsOutputGraylogStreamDelete,
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
						"GLEF",
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
		},
	}
}

func resourceDbaasLogsOutputGraylogStreamImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not service_name/id formatted")
	}
	serviceName := splitId[0]
	id := splitId[1]
	d.SetId(id)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceDbaasLogsOutputGraylogStreamCreate(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error calling post %s:\n\t %q", endpoint, err)
	}

	// Wait for operation status
	op, err := waitForDbaasLogsOperation(config.OVHClient, serviceName, res.OperationId)
	if err != nil {
		return err
	}

	id := op.StreamId
	if id == nil {
		return fmt.Errorf("Stream Id is nil. This should not happen: operation is %s/%s", serviceName, res.OperationId)
	}

	d.SetId(*id)

	return resourceDbaasLogsOutputGraylogStreamRead(d, meta)
}

func resourceDbaasLogsOutputGraylogStreamUpdate(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error calling Put %s:\n\t %q", endpoint, err)
	}

	// Wait for operation status
	if _, err := waitForDbaasLogsOperation(config.OVHClient, serviceName, res.OperationId); err != nil {
		return err
	}

	return resourceDbaasLogsOutputGraylogStreamRead(d, meta)
}

func resourceDbaasLogsOutputGraylogStreamRead(d *schema.ResourceData, meta interface{}) error {
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
		return helpers.CheckDeleted(d, err, endpoint)
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

func resourceDbaasLogsOutputGraylogStreamDelete(d *schema.ResourceData, meta interface{}) error {
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
		return helpers.CheckDeleted(d, err, endpoint)
	}

	// Wait for operation status
	if _, err := waitForDbaasLogsOperation(config.OVHClient, serviceName, res.OperationId); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

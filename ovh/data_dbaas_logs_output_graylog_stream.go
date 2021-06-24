package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDbaasLogsOutputGraylogStream() *schema.Resource {
	return &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error {
			return dataSourceDbaasLogsOutputGraylogStreamRead(d, meta)
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The service name",
				Required:    true,
			},
			"title": {
				Type:        schema.TypeString,
				Description: "Stream description",
				Required:    true,
			},

			// computed
			"can_alert": {
				Type:        schema.TypeBool,
				Description: "Indicates if the current user can create alert on the stream",
				Computed:    true,
			},
			"cold_storage_compression": {
				Type:        schema.TypeString,
				Description: "Cold storage compression method",
				Computed:    true,
			},
			"cold_storage_content": {
				Type:        schema.TypeString,
				Description: "ColdStorage content",
				Computed:    true,
			},
			"cold_storage_enabled": {
				Type:        schema.TypeBool,
				Description: "Is Cold storage enabled?",
				Computed:    true,
			},
			"cold_storage_notify_enabled": {
				Type:        schema.TypeBool,
				Description: "Notify on new Cold storage archive",
				Computed:    true,
			},
			"cold_storage_retention": {
				Type:        schema.TypeInt,
				Description: "Cold storage retention in year",
				Computed:    true,
			},
			"cold_storage_target": {
				Type:        schema.TypeString,
				Description: "ColdStorage destination",
				Computed:    true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Description: "Stream creation",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Stream description",
				Computed:    true,
			},
			"indexing_enabled": {
				Type:        schema.TypeBool,
				Description: "Enable ES indexing",
				Computed:    true,
			},
			"indexing_max_size": {
				Type:        schema.TypeInt,
				Description: "Maximum indexing size (in GB)",
				Computed:    true,
			},
			"indexing_notify_enabled": {
				Type:        schema.TypeBool,
				Description: "If set, notify when size is near 80, 90 or 100 % of the maximum configured setting",
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
			"parent_stream_id": {
				Type:        schema.TypeString,
				Description: "Parent stream ID",
				Computed:    true,
			},
			"pause_indexing_on_max_size": {
				Type:        schema.TypeBool,
				Description: "If set, pause indexing when maximum size is reach",
				Computed:    true,
			},
			"retention_id": {
				Type:        schema.TypeString,
				Description: "Retention ID",
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
			"web_socket_enabled": {
				Type:        schema.TypeString,
				Description: "Enable Websocket",
				Computed:    true,
			},
		},
	}
}

func dataSourceDbaasLogsOutputGraylogStreamRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	titleFilter := d.Get("title").(string)

	log.Printf("[DEBUG] Will read dbaas logs output graylog streams")
	res := []string{}
	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/output/graylog/stream",
		url.PathEscape(serviceName),
	)
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
	}
	streams := []*DbaasLogsOutputGraylogStream{}

	for _, id := range res {
		log.Printf("[DEBUG] Will read dbaas logs output graylog stream id : %s/%s", serviceName, id)

		endpoint := fmt.Sprintf(
			"/dbaas/logs/%s/output/graylog/stream/%s",
			url.PathEscape(serviceName),
			url.PathEscape(id),
		)

		stream := &DbaasLogsOutputGraylogStream{}
		if err := config.OVHClient.Get(endpoint, &stream); err != nil {
			return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
		}

		log.Printf("[INFO]Comparing : %s ? %s",
			strings.ToLower(stream.Title),
			strings.ToLower(titleFilter),
		)

		if strings.ToLower(stream.Title) == strings.ToLower(titleFilter) {
			streams = append(streams, stream)
		}
	}

	if len(streams) == 0 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	if len(streams) > 1 {
		return fmt.Errorf("Your query returned more than one result. " +
			"Please change your search criteria and try again.")
	}

	for k, v := range streams[0].ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	return nil
}

package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceDbaasLogsOutputGraylogStreamURL returns the list of URLs for a Graylog stream.
func dataSourceDbaasLogsOutputGraylogStreamURL() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDbaasLogsOutputGraylogStreamURLRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The service name",
				Required:    true,
			},
			"stream_id": {
				Type:        schema.TypeString,
				Description: "Stream ID",
				Required:    true,
			},
			"url": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL address",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL type",
						},
					},
				},
			},
		},
	}
}

func dataSourceDbaasLogsOutputGraylogStreamURLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	streamID := d.Get("stream_id").(string)

	log.Printf("[DEBUG] Will read URLs for dbaas logs output graylog stream: %s/%s", serviceName, streamID)
	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/output/graylog/stream/%s/url",
		url.PathEscape(serviceName),
		url.PathEscape(streamID),
	)

	var urls []DbaasLogsOutputGraylogStreamURL
	if err := config.OVHClient.Get(endpoint, &urls); err != nil {
		return diag.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
	}

	list := make([]map[string]interface{}, len(urls))
	for i, u := range urls {
		m := map[string]interface{}{
			"address": u.Address,
			"type":    u.Type,
		}
		list[i] = m
	}

	if err := d.Set("url", list); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", serviceName, streamID))
	return nil
}

package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDbaasLogsOutputOpensearchIndex() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDbaasLogsOutputOpensearchIndexRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The service name",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Index name",
				Required:    true,
			},

			// computed
			"alert_notify_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set, notify when size is near 80, 90 or 100 % of its maximum capacity",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operation creation",
			},
			"current_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Current Index size (in bytes)",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Index description",
			},
			"index_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Index ID",
			},
			"is_editable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if you are allowed to edit entry",
			},
			"max_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum index size (in bytes)",
			},
			"nb_shard": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of shard",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operation last update",
			},
		},
	}
}

func dataSourceDbaasLogsOutputOpensearchIndexRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	nameFilter := d.Get("name").(string)

	log.Printf("[DEBUG] Will read dbaas logs output opensearch indexes")
	res := []string{}
	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/output/opensearch/index",
		url.PathEscape(serviceName),
	)
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return diag.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
	}
	indexes := []*DbaasLogsOutputOpensearchIndex{}

	for _, id := range res {
		log.Printf("[DEBUG] Will read dbaas logs output opensearch index id: %s/%s", serviceName, id)

		endpoint := fmt.Sprintf(
			"/dbaas/logs/%s/output/opensearch/index/%s",
			url.PathEscape(serviceName),
			url.PathEscape(id),
		)

		index := &DbaasLogsOutputOpensearchIndex{}
		if err := config.OVHClient.Get(endpoint, &index); err != nil {
			return diag.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
		}

		log.Printf("[DEBUG]Comparing: %s ? %s",
			strings.ToLower(index.Name),
			strings.ToLower(nameFilter),
		)

		if strings.EqualFold(index.Name, nameFilter) {
			indexes = append(indexes, index)
		}
	}

	if len(indexes) == 0 {
		return diag.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}
	if len(indexes) > 1 {
		return diag.Errorf("Your query returned more than one result. Please change your search criteria and try again.")
	}

	for k, v := range indexes[0].ToMap() {
		if k != "index_id" {
			d.Set(k, v)
		} else {
			d.Set("index_id", fmt.Sprint(v))
			d.SetId(fmt.Sprint(v))
		}
	}
	return nil
}

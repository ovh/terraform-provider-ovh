package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func dataSourceCloudProjectDatabaseIntegration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProjectDatabaseIntegrationRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"engine": {
				Type:        schema.TypeString,
				Description: "Name of the engine of the service",
				Required:    true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value == "mongodb" {
						errors = append(errors, fmt.Errorf("value %s is not a valid engine for integration", value))
					}
					return
				},
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Cluster ID",
				Required:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "Integration ID",
				Required:    true,
			},

			//Computed
			"destination_service_id": {
				Type:        schema.TypeString,
				Description: "ID of the destination service",
				Computed:    true,
			},
			"parameters": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Parameters for the integration",
				Computed:    true,
			},
			"source_service_id": {
				Type:        schema.TypeString,
				Description: "ID of the source service",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Current status of the integration",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Type of the integration",
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudProjectDatabaseIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterID := d.Get("cluster_id").(string)
	id := d.Get("id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/integration/%s",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterID),
		url.PathEscape(id),
	)

	res := &CloudProjectDatabaseIntegrationResponse{}

	log.Printf("[DEBUG] Will read acl %s from cluster %s from project %s", id, clusterID, serviceName)
	if err := config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	log.Printf("[DEBUG] Read integration %+v", res)
	return nil
}

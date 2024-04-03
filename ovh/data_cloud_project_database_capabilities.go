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

func dataSourceCloudProjectDatabaseCapabilities() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProjectDatabaseCapabilitiesRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},

			// Computed
			"engines": {
				Type:        schema.TypeSet,
				Description: "Database engines available",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_version": {
							Type:        schema.TypeString,
							Description: "Default version used for the engine",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Description of the engine",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Engine name",
							Computed:    true,
						},
						"ssl_modes": {
							Type:        schema.TypeSet,
							Description: "SSL modes for this engine",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"versions": {
							Type:        schema.TypeSet,
							Description: "Versions available for this engine",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"flavors": {
				Type:        schema.TypeSet,
				Description: "Flavors available",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"core": {
							Type:        schema.TypeInt,
							Description: "Flavor core number",
							Computed:    true,
						},
						"memory": {
							Type:        schema.TypeInt,
							Description: "Flavor ram size in GB",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the flavor",
							Computed:    true,
						},
						"storage": {
							Type:        schema.TypeInt,
							Description: "Flavor disk size in GB",
							Computed:    true,
						},
					},
				},
			},
			"options": {
				Type:        schema.TypeSet,
				Description: "Options available",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the option",
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Type of the option",
							Computed:    true,
						},
					},
				},
			},
			"plans": {
				Type:        schema.TypeSet,
				Description: "Plans available",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_retention": {
							Type:        schema.TypeString,
							Description: "Automatic backup retention duration",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Description of the plan",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the plan",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudProjectDatabaseCapabilitiesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	capabilitiesEndpoint := fmt.Sprintf("/cloud/project/%s/database/capabilities",
		url.PathEscape(serviceName),
	)
	capabilitiesRes := &CloudProjectDatabaseCapabilitiesResponse{}

	log.Printf("[DEBUG] Will read capabilities from project %s", serviceName)
	if err := config.OVHClient.GetWithContext(ctx, capabilitiesEndpoint, capabilitiesRes); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, capabilitiesEndpoint))
	}

	d.SetId(serviceName)
	for k, v := range capabilitiesRes.ToMap() {
		d.Set(k, v)
	}

	log.Printf("[DEBUG] Read capabilities %+v", capabilitiesRes)
	return nil
}

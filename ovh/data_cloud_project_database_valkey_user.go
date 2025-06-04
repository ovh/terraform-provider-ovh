package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudProjectDatabaseValkeyUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProjectDatabaseValkeyUserRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Id of the database cluster",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the user",
				Required:    true,
			},

			//Computed
			"categories": {
				Type:        schema.TypeSet,
				Description: "Categories of the user",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"channels": {
				Type:        schema.TypeSet,
				Description: "Channels of the user",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"commands": {
				Type:        schema.TypeSet,
				Description: "Commands of the user",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Type:        schema.TypeString,
				Description: "Date of the creation of the user",
				Computed:    true,
			},
			"keys": {
				Type:        schema.TypeSet,
				Description: "Keys of the user",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Current status of the user",
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudProjectDatabaseValkeyUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterID := d.Get("cluster_id").(string)

	listEndpoint := fmt.Sprintf("/cloud/project/%s/database/valkey/%s/user",
		url.PathEscape(serviceName),
		url.PathEscape(clusterID),
	)

	listRes := make([]string, 0)

	log.Printf("[DEBUG] Will read users from cluster %s from project %s", clusterID, serviceName)
	if err := config.OVHClient.GetWithContext(ctx, listEndpoint, &listRes); err != nil {
		return diag.Errorf("Error calling GET %s:\n\t %q", listEndpoint, err)
	}

	name := d.Get("name").(string)
	for _, id := range listRes {
		endpoint := fmt.Sprintf("/cloud/project/%s/database/valkey/%s/user/%s",
			url.PathEscape(serviceName),
			url.PathEscape(clusterID),
			url.PathEscape(id),
		)
		res := &CloudProjectDatabaseRedisUserResponse{}

		log.Printf("[DEBUG] Will read user %s from cluster %s from project %s", id, clusterID, serviceName)
		if err := config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
			return diag.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
		}

		if res.Username == name {
			for k, v := range res.ToMap() {
				if k != "id" {
					d.Set(k, v)
				} else {
					d.SetId(fmt.Sprint(v))
				}
			}
			log.Printf("[DEBUG] Read user %+v", res)
			return nil
		}
	}

	return diag.Errorf("User name %s not found for cluster %s from project %s", name, clusterID, serviceName)
}

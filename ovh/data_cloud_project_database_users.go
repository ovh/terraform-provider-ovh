package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

func dataSourceCloudProjectDatabaseUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProjectDatabaseUsersRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"engine": {
				Type:             schema.TypeString,
				Description:      "Name of the engine of the service",
				Required:         true,
				ValidateDiagFunc: helpers.ValidateDiagEnum(engines),
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Cluster ID",
				Required:    true,
			},

			//Computed
			"user_ids": {
				Type:        schema.TypeList,
				Description: "List of users ids",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"users": {
				Type:        schema.TypeList,
				Description: "List of users with their details",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "ID of the user",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the user",
							Computed:    true,
						},
						"created_at": {
							Type:        schema.TypeString,
							Description: "Date of the creation of the user",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Current status of the user",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudProjectDatabaseUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterID := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/user",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterID),
	)

	res := make([]string, 0)

	log.Printf("[DEBUG] Will read users from cluster %s from project %s", clusterID, serviceName)
	if err := config.OVHClient.GetWithContext(ctx, endpoint, &res); err != nil {
		return diag.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(res)

	users := make([]map[string]interface{}, 0, len(res))
	for _, id := range res {
		userEndpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/user/%s",
			url.PathEscape(serviceName),
			url.PathEscape(engine),
			url.PathEscape(clusterID),
			url.PathEscape(id),
		)
		user := &CloudProjectDatabaseUserResponse{}

		log.Printf("[DEBUG] Will read user %s from cluster %s from project %s", id, clusterID, serviceName)
		if err := config.OVHClient.GetWithContext(ctx, userEndpoint, user); err != nil {
			return diag.Errorf("Error calling GET %s:\n\t %q", userEndpoint, err)
		}

		users = append(users, user.ToMap())
	}

	d.SetId(hashcode.Strings(res))
	d.Set("user_ids", res)
	d.Set("users", users)

	return nil
}

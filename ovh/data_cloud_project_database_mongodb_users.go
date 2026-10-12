package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

func dataSourceCloudProjectDatabaseMongodbUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProjectDatabaseMongodbUsersRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Cluster ID",
				Required:    true,
			},

			//Computed
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
							Description: "Name of the user with the authentication database in the format name@authDB",
							Computed:    true,
						},
						"created_at": {
							Type:        schema.TypeString,
							Description: "Date of the creation of the user",
							Computed:    true,
						},
						"roles": {
							Type:        schema.TypeSet,
							Description: "Roles the user belongs to",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
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

func dataSourceCloudProjectDatabaseMongodbUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterID := d.Get("cluster_id").(string)

	listEndpoint := fmt.Sprintf("/cloud/project/%s/database/mongodb/%s/user",
		url.PathEscape(serviceName),
		url.PathEscape(clusterID),
	)

	listRes := make([]string, 0)

	log.Printf("[DEBUG] Will read users from cluster %s from project %s", clusterID, serviceName)
	if err := config.OVHClient.GetWithContext(ctx, listEndpoint, &listRes); err != nil {
		return diag.Errorf("Error calling GET %s:\n\t %q", listEndpoint, err)
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(listRes)

	users := make([]map[string]interface{}, 0, len(listRes))
	for _, id := range listRes {
		endpoint := fmt.Sprintf("/cloud/project/%s/database/mongodb/%s/user/%s",
			url.PathEscape(serviceName),
			url.PathEscape(clusterID),
			url.PathEscape(id),
		)
		res := &CloudProjectDatabaseMongodbUserResponse{}

		log.Printf("[DEBUG] Will read user %s from cluster %s from project %s", id, clusterID, serviceName)
		if err := config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
			return diag.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
		}

		users = append(users, res.toMap())
	}

	d.SetId(hashcode.Strings(listRes))
	d.Set("users", users)

	return nil
}

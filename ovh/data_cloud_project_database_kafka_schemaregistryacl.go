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

func dataSourceCloudProjectDatabaseKafkaSchemaRegistryAcl() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProjectDatabaseKafkaSchemaRegistryAclRead,
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
			"id": {
				Type:        schema.TypeString,
				Description: "Shema registry ACL ID",
				Required:    true,
			},

			// Computed
			"permission": {
				Type:        schema.TypeString,
				Description: "Permission to give to this username on this resource",
				Computed:    true,
			},
			"resource": {
				Type:        schema.TypeString,
				Description: "Resource affected by this ACL",
				Computed:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "Username affected by this ACL",
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudProjectDatabaseKafkaSchemaRegistryAclRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterID := d.Get("cluster_id").(string)
	id := d.Get("id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/schemaRegistryAcl/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterID),
		url.PathEscape(id),
	)
	res := &CloudProjectDatabaseKafkaAclResponse{}

	log.Printf("[DEBUG] Will read schema registry ACL %s from cluster %s from project %s", id, clusterID, serviceName)
	if err := config.OVHClient.Get(endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	log.Printf("[DEBUG] Read ACL %+v", res)
	return nil
}

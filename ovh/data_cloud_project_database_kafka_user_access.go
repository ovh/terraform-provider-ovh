package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceCloudProjectDatabaseKafkaUserAccess() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProjectDatabaseKafkaUserAccessRead,
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
			"user_id": {
				Type:        schema.TypeString,
				Description: "Id of the user",
				Required:    true,
			},

			// Computed
			"cert": {
				Type:        schema.TypeString,
				Description: "User cert",
				Computed:    true,
			},
			"key": {
				Type:        schema.TypeString,
				Description: "User key for the cert",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func dataSourceCloudProjectDatabaseKafkaUserAccessRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	userId := d.Get("user_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/user/%s/access",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(userId),
	)
	res := &CloudProjectDatabaseKafkaUserAccessResponse{}

	log.Printf("[DEBUG] Will read certificates of user %s from cluster %s from project %s", userId, clusterId, serviceName)
	if err := config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	d.SetId(strconv.Itoa(hashcode.String(res.Cert)))
	for k, v := range res.ToMap() {
		d.Set(k, v)
	}

	log.Printf("[DEBUG] Read certificates %+v", res)
	return nil
}

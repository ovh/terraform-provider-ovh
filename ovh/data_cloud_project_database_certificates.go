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

func dataSourceCloudProjectDatabaseCertificates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProjectDatabaseCertificatesRead,
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
				ValidateDiagFunc: helpers.ValidateDiagEnum([]string{"cassandra", "kafka", "mysql", "postgresql"}),
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Id of the database cluster",
				Required:    true,
			},

			// Computed
			"ca": {
				Type:        schema.TypeString,
				Description: "CA certificate used for the service",
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudProjectDatabaseCertificatesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterID := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/certificates",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterID),
	)
	res := &CloudProjectDatabaseCertificatesResponse{}

	log.Printf("[DEBUG] Will read certificates from cluster %s from project %s", clusterID, serviceName)
	if err := config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	d.SetId(strconv.Itoa(hashcode.String(res.Ca)))
	for k, v := range res.ToMap() {
		d.Set(k, v)
	}

	log.Printf("[DEBUG] Read certificates %+v", res)
	return nil
}

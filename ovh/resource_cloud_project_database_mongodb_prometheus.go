package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceCloudProjectDatabaseMongodbPrometheus() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudProjectDatabaseMongodbPrometheusCreate,
		ReadContext:   resourceCloudProjectDatabaseMongodbPrometheusRead,
		DeleteContext: resourceCloudProjectDatabaseMongodbPrometheusDelete,
		UpdateContext: resourceCloudProjectDatabaseMongodbPrometheusUpdate,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectDatabaseMongodbPrometheusImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Id of the database cluster",
				ForceNew:    true,
				Required:    true,
			},
			"password_reset": {
				Type:        schema.TypeString,
				Description: "Arbitrary string to change to trigger a password update",
				Optional:    true,
			},

			//Computed
			"username": {
				Type:        schema.TypeString,
				Description: "Name of the user",
				Computed:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "Password of the user",
				Sensitive:   true,
				Computed:    true,
			},
			"srv_domain": {
				Type:        schema.TypeString,
				Description: "Name of the srv domain endpoint",
				Computed:    true,
			},
		},
	}
}

func resourceCloudProjectDatabaseMongodbPrometheusImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenID := d.Id()
	n := 2
	splitID := strings.SplitN(givenID, "/", n)
	if len(splitID) != n {
		return nil, fmt.Errorf("import Id is not service_name/cluster_id formatted")
	}
	serviceName := splitID[0]
	id := splitID[1]
	d.SetId(id)
	d.Set("cluster_id", id)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectDatabaseMongodbPrometheusCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return enableCloudProjectDatabasePrometheus(ctx, d, meta, "mongodb", true, resourceCloudProjectDatabaseMongodbPrometheusUpdate)
}

func resourceCloudProjectDatabaseMongodbPrometheusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/mongodb/%s/prometheus",
		url.PathEscape(serviceName),
		url.PathEscape(d.Id()),
	)
	res := &CloudProjectDatabaseMongodbPrometheusEndpointResponse{}

	log.Printf("[DEBUG] Will read prometheus of database %s from project: %s", d.Id(), serviceName)
	if err := config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	for k, v := range res.toMap() {
		d.Set(k, v)
	}

	return nil
}

func resourceCloudProjectDatabaseMongodbPrometheusUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return updateCloudProjectDatabasePrometheus(ctx, d, meta, "mongodb", resourceCloudProjectDatabaseMongodbPrometheusRead)
}

func resourceCloudProjectDatabaseMongodbPrometheusDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return enableCloudProjectDatabasePrometheus(ctx, d, meta, "mongodb", false, nil)
}

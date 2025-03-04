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
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceCloudProjectDatabasePrometheus() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudProjectDatabasePrometheusCreate,
		ReadContext:   resourceCloudProjectDatabasePrometheusRead,
		DeleteContext: resourceCloudProjectDatabasePrometheusDelete,
		UpdateContext: resourceCloudProjectDatabasePrometheusUpdate,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectDatabasePrometheusImportState,
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
			"engine": {
				Type:             schema.TypeString,
				Description:      "Name of the engine of the service",
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: helpers.ValidateDiagEnum([]string{"cassandra", "kafka", "kafkaConnect", "kafkaMirrorMaker", "mysql", "opensearch", "postgresql", "redis"}),
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
			"targets": {
				Type:        schema.TypeList,
				Description: "List of all endpoint targets",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Description: "Host of the endpoint",
							Computed:    true,
						},
						"port": {
							Type:        schema.TypeInt,
							Description: "Connection port for the endpoint",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceCloudProjectDatabasePrometheusImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenID := d.Id()
	n := 3
	splitID := strings.SplitN(givenID, "/", n)
	if len(splitID) != n {
		return nil, fmt.Errorf("import Id is not service_name/engine/cluster_id formatted")
	}
	serviceName := splitID[0]
	engine := splitID[1]
	id := splitID[2]
	d.SetId(id)
	d.Set("cluster_id", id)
	d.Set("engine", engine)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectDatabasePrometheusCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	engine := d.Get("engine").(string)
	return enableCloudProjectDatabasePrometheus(ctx, d, meta, engine, true, resourceCloudProjectDatabasePrometheusUpdate)
}

func resourceCloudProjectDatabasePrometheusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/prometheus",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(d.Id()),
	)
	res := &CloudProjectDatabasePrometheusEndpointResponse{}

	log.Printf("[DEBUG] Will read prometheus of cluster %s from project: %s", d.Id(), serviceName)
	if err := config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	for k, v := range res.toMap() {
		d.Set(k, v)
	}

	return nil
}

func resourceCloudProjectDatabasePrometheusUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	engine := d.Get("engine").(string)
	return updateCloudProjectDatabasePrometheus(ctx, d, meta, engine, resourceCloudProjectDatabasePrometheusRead)
}

func resourceCloudProjectDatabasePrometheusDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	engine := d.Get("engine").(string)
	return enableCloudProjectDatabasePrometheus(ctx, d, meta, engine, false, nil)
}

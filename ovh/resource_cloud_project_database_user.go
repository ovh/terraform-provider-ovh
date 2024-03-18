package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceCloudProjectDatabaseUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudProjectDatabaseUserCreate,
		ReadContext:   resourceCloudProjectDatabaseUserRead,
		UpdateContext: resourceCloudProjectDatabaseUserUpdate,
		DeleteContext: resourceCloudProjectDatabaseUserDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectDatabaseUserImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
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
				Type:         schema.TypeString,
				Description:  "Name of the engine of the service",
				ForceNew:     true,
				Required:     true,
				ValidateFunc: helpers.ValidateEnum([]string{"cassandra", "mysql", "kafka", "kafkaConnect", "grafana"}),
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Id of the database cluster",
				ForceNew:    true,
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the user",
				ForceNew:    true,
				Required:    true,
			},
			"password_reset": {
				Type:        schema.TypeString,
				Description: "Arbitrary string to change to trigger a password update",
				Optional:    true,
			},

			//Computed
			"created_at": {
				Type:        schema.TypeString,
				Description: "Date of the creation of the user",
				Computed:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "Password of the user",
				Sensitive:   true,
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Current status of the user",
				Computed:    true,
			},
		},

		CustomizeDiff: customdiff.ComputedIf(
			"password",
			func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
				return d.HasChange("password_reset")
			},
		),
	}
}

func resourceCloudProjectDatabaseUserImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	n := 4
	splitId := strings.SplitN(givenId, "/", n)
	if len(splitId) != n {
		return nil, fmt.Errorf("import Id is not service_name/engine/cluster_id/id formatted")
	}
	serviceName := splitId[0]
	engine := splitId[1]
	clusterId := splitId[2]
	id := splitId[3]
	d.SetId(id)
	d.Set("cluster_id", clusterId)
	d.Set("engine", engine)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectDatabaseUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	engine := d.Get("engine").(string)
	f := func() interface{} {
		return (&CloudProjectDatabaseUserCreateOpts{}).FromResource(d)
	}
	return postCloudProjectDatabaseUser(ctx, d, meta, engine, dataSourceCloudProjectDatabaseUserRead, resourceCloudProjectDatabaseUserRead, resourceCloudProjectDatabaseUserUpdate, f)
}

func resourceCloudProjectDatabaseUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/user/%s",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)

	res := &CloudProjectDatabaseUserResponse{}

	log.Printf("[DEBUG] Will read user %s from cluster %s from project %s", id, clusterId, serviceName)
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

	return nil
}

func resourceCloudProjectDatabaseUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/user/%s/credentials/reset",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)
	res := &CloudProjectDatabaseUserResponse{}
	log.Printf("[DEBUG] Will update user password for cluster %s from project %s", clusterId, serviceName)
	err := postFuncCloudProjectDatabaseUser(ctx, d, meta, engine, endpoint, nil, res, schema.TimeoutUpdate)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCloudProjectDatabaseUserRead(ctx, d, meta)
}

func resourceCloudProjectDatabaseUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	engine := d.Get("engine").(string)
	return deleteCloudProjectDatabaseUser(ctx, d, meta, engine)
}

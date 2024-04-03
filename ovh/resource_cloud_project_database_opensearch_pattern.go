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

func resourceCloudProjectDatabaseOpensearchPattern() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudProjectDatabaseOpensearchPatternCreate,
		ReadContext:   resourceCloudProjectDatabaseOpensearchPatternRead,
		DeleteContext: resourceCloudProjectDatabaseOpensearchPatternDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectDatabaseOpensearchPatternImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
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
			"max_index_count": {
				Type:        schema.TypeInt,
				Description: "Maximum number of index for this pattern",
				ForceNew:    true,
				Optional:    true,
			},
			"pattern": {
				Type:        schema.TypeString,
				Description: "Pattern format",
				ForceNew:    true,
				Required:    true,
			},
		},
	}
}

func resourceCloudProjectDatabaseOpensearchPatternImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	n := 3
	splitId := strings.SplitN(givenId, "/", n)
	if len(splitId) != n {
		return nil, fmt.Errorf("import Id is not service_name/cluster_id/id formatted")
	}
	serviceName := splitId[0]
	clusterId := splitId[1]
	id := splitId[2]
	d.SetId(id)
	d.Set("cluster_id", clusterId)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectDatabaseOpensearchPatternCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/opensearch/%s/pattern",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
	)
	params := (&CloudProjectDatabaseOpensearchPatternCreateOpts{}).FromResource(d)
	res := &CloudProjectDatabaseOpensearchPatternResponse{}

	log.Printf("[DEBUG] Will create pattern: %+v for cluster %s from project %s", params, clusterId, serviceName)
	err := config.OVHClient.PostWithContext(ctx, endpoint, params, res)
	if err != nil {
		return diag.Errorf("calling Post %s with params %+v:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for topic %s to be READY", res.Id)
	err = waitForCloudProjectDatabaseOpensearchPatternReady(ctx, config.OVHClient, serviceName, clusterId, res.Id, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("timeout while waiting topic %s to be READY: %s", res.Id, err.Error())
	}
	log.Printf("[DEBUG] topic %s is READY", res.Id)

	d.SetId(res.Id)

	return resourceCloudProjectDatabaseOpensearchPatternRead(ctx, d, meta)
}

func resourceCloudProjectDatabaseOpensearchPatternRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/opensearch/%s/pattern/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)
	res := &CloudProjectDatabaseOpensearchPatternResponse{}

	log.Printf("[DEBUG] Will read pattern %s from cluster %s from project %s", id, clusterId, serviceName)
	if err := config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	log.Printf("[DEBUG] Read pattern %+v", res)
	return nil
}

func resourceCloudProjectDatabaseOpensearchPatternDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/opensearch/%s/pattern/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)

	log.Printf("[DEBUG] Will delete pattern %s from cluster %s from project %s", id, clusterId, serviceName)
	err := config.OVHClient.DeleteWithContext(ctx, endpoint, nil)
	if err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	log.Printf("[DEBUG] Waiting for pattern %s to be DELETED", id)
	err = waitForCloudProjectDatabaseOpensearchPatternDeleted(ctx, config.OVHClient, serviceName, clusterId, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("timeout while waiting pattern %s to be DELETED: %s", id, err.Error())
	}
	log.Printf("[DEBUG] pattern %s is DELETED", id)

	d.SetId("")

	return nil
}

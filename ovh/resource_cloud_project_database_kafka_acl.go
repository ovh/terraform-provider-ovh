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

func resourceCloudProjectDatabaseKafkaACL() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudProjectDatabaseKafkaACLCreate,
		ReadContext:   resourceCloudProjectDatabaseKafkaACLRead,
		DeleteContext: resourceCloudProjectDatabaseKafkaACLDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectDatabaseKafkaACLImportState,
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
			"permission": {
				Type:        schema.TypeString,
				Description: "Permission to give to this username on this topic",
				ForceNew:    true,
				Required:    true,
			},
			"topic": {
				Type:        schema.TypeString,
				Description: "Topic affected by this acl",
				ForceNew:    true,
				Required:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "Username affected by this acl",
				ForceNew:    true,
				Required:    true,
			},
		},
	}
}

func resourceCloudProjectDatabaseKafkaACLImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenID := d.Id()
	n := 3
	splitID := strings.SplitN(givenID, "/", n)
	if len(splitID) != n {
		return nil, fmt.Errorf("import Id is not service_name/cluster_id/id formatted")
	}
	serviceName := splitID[0]
	clusterID := splitID[1]
	id := splitID[2]
	d.SetId(id)
	d.Set("cluster_id", clusterID)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectDatabaseKafkaACLCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterID := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/acl",
		url.PathEscape(serviceName),
		url.PathEscape(clusterID),
	)
	params := (&CloudProjectDatabaseKafkaACLCreateOpts{}).FromResource(d)
	res := &CloudProjectDatabaseKafkaACLResponse{}

	log.Printf("[DEBUG] Will create acl: %+v for cluster %s from project %s", params, clusterID, serviceName)
	err := config.OVHClient.PostWithContext(ctx, endpoint, params, res)
	if err != nil {
		return diag.Errorf("calling Post %s with params %+v:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for acl %s to be READY", res.ID)
	err = waitForCloudProjectDatabaseKafkaACLReady(ctx, config.OVHClient, serviceName, clusterID, res.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("timeout while waiting acl %s to be READY: %s", res.ID, err.Error())
	}
	log.Printf("[DEBUG] acl %s is READY", res.ID)

	d.SetId(res.ID)

	return resourceCloudProjectDatabaseKafkaACLRead(ctx, d, meta)
}

func resourceCloudProjectDatabaseKafkaACLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterID := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/acl/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterID),
		url.PathEscape(id),
	)
	res := &CloudProjectDatabaseKafkaACLResponse{}

	log.Printf("[DEBUG] Will read acl %s from cluster %s from project %s", id, clusterID, serviceName)
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

	log.Printf("[DEBUG] Read acl %+v", res)
	return nil
}

func resourceCloudProjectDatabaseKafkaACLDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterID := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/acl/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterID),
		url.PathEscape(id),
	)

	log.Printf("[DEBUG] Will delete acl  %s from cluster %s from project %s", id, clusterID, serviceName)
	err := config.OVHClient.DeleteWithContext(ctx, endpoint, nil)
	if err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	log.Printf("[DEBUG] Waiting for acl %s to be DELETED", id)
	err = waitForCloudProjectDatabaseKafkaACLDeleted(ctx, config.OVHClient, serviceName, clusterID, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("timeout while waiting acl %s to be DELETED: %s", id, err.Error())
	}
	log.Printf("[DEBUG] acl %s is DELETED", id)

	d.SetId("")

	return nil
}

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

func resourceCloudProjectDatabaseKafkaAcl() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudProjectDatabaseKafkaAclCreate,
		ReadContext:   resourceCloudProjectDatabaseKafkaAclRead,
		DeleteContext: resourceCloudProjectDatabaseKafkaAclDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectDatabaseKafkaAclImportState,
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

func resourceCloudProjectDatabaseKafkaAclImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

func resourceCloudProjectDatabaseKafkaAclCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/acl",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
	)
	params := (&CloudProjectDatabaseKafkaAclCreateOpts{}).FromResource(d)
	res := &CloudProjectDatabaseKafkaAclResponse{}

	log.Printf("[DEBUG] Will create acl: %+v for cluster %s from project %s", params, clusterId, serviceName)
	err := config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		return diag.Errorf("calling Post %s with params %+v:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for acl %s to be READY", res.Id)
	err = waitForCloudProjectDatabaseKafkaAclReady(ctx, config.OVHClient, serviceName, clusterId, res.Id, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("timeout while waiting acl %s to be READY: %s", res.Id, err.Error())
	}
	log.Printf("[DEBUG] acl %s is READY", res.Id)

	d.SetId(res.Id)

	return resourceCloudProjectDatabaseKafkaAclRead(ctx, d, meta)
}

func resourceCloudProjectDatabaseKafkaAclRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/acl/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)
	res := &CloudProjectDatabaseKafkaAclResponse{}

	log.Printf("[DEBUG] Will read acl %s from cluster %s from project %s", id, clusterId, serviceName)
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

	log.Printf("[DEBUG] Read acl %+v", res)
	return nil
}

func resourceCloudProjectDatabaseKafkaAclDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/acl/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)

	log.Printf("[DEBUG] Will delete acl  %s from cluster %s from project %s", id, clusterId, serviceName)
	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	log.Printf("[DEBUG] Waiting for acl %s to be DELETED", id)
	err = waitForCloudProjectDatabaseKafkaAclDeleted(ctx, config.OVHClient, serviceName, clusterId, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("timeout while waiting acl %s to be DELETED: %s", id, err.Error())
	}
	log.Printf("[DEBUG] acl %s is DELETED", id)

	d.SetId("")

	return nil
}

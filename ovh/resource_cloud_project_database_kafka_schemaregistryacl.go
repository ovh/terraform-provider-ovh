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

func resourceCloudProjectDatabaseKafkaSchemaRegistryAcl() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudProjectDatabaseKafkaSchemaRegistryAclCreate,
		ReadContext:   resourceCloudProjectDatabaseKafkaSchemaRegistryAclRead,
		DeleteContext: resourceCloudProjectDatabaseKafkaSchemaRegistryAclDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectDatabaseKafkaSchemaRegistryAclImportState,
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
				Description: "Permission to give to this username on this resource",
				ForceNew:    true,
				Required:    true,
			},
			"resource": {
				Type:        schema.TypeString,
				Description: "Resource affected by this acl",
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

func resourceCloudProjectDatabaseKafkaSchemaRegistryAclImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	n := 3
	splitId := strings.SplitN(givenId, "/", n)
	if len(splitId) != n {
		return nil, fmt.Errorf("Import Id is not service_name/cluster_id/id formatted")
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

func resourceCloudProjectDatabaseKafkaSchemaRegistryAclCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/schemaRegistryAcl",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
	)
	params := (&CloudProjectDatabaseKafkaSchemaRegistryAclCreateOpts{}).FromResource(d)
	res := &CloudProjectDatabaseKafkaSchemaRegistryAclResponse{}

	log.Printf("[DEBUG] Will create schema registry acl: %+v for cluster %s from project %s", params, clusterId, serviceName)
	err := config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		return diag.Errorf("calling Post %s with params %+v:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for schema registry acl %s to be READY", res.Id)
	err = waitForCloudProjectDatabaseKafkaSchemaRegistryAclReady(ctx, config.OVHClient, serviceName, clusterId, res.Id, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("timeout while waiting schema registry ACL %s to be READY: %s", res.Id, err.Error())
	}
	log.Printf("[DEBUG] schema registry acl %s is READY", res.Id)

	d.SetId(res.Id)

	return resourceCloudProjectDatabaseKafkaSchemaRegistryAclRead(ctx, d, meta)
}

func resourceCloudProjectDatabaseKafkaSchemaRegistryAclRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/schemaRegistryAcl/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)
	res := &CloudProjectDatabaseKafkaSchemaRegistryAclResponse{}

	log.Printf("[DEBUG] Will read schema registry acl %s from cluster %s from project %s", id, clusterId, serviceName)
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

	log.Printf("[DEBUG] Read schema registry ACL %+v", res)
	return nil
}

func resourceCloudProjectDatabaseKafkaSchemaRegistryAclDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/schemaRegistryAcl/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)

	log.Printf("[DEBUG] Will delete schema registry acl  %s from cluster %s from project %s", id, clusterId, serviceName)
	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	log.Printf("[DEBUG] Waiting for schema registry acl %s to be DELETED", id)
	err = waitForCloudProjectDatabaseKafkaSchemaRegistryAclDeleted(ctx, config.OVHClient, serviceName, clusterId, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("timeout while waiting schema registry ACL %s to be DELETED: %s", id, err.Error())
	}
	log.Printf("[DEBUG] schema registry acl %s is DELETED", id)

	d.SetId("")

	return nil
}

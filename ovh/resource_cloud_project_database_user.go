package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceCloudProjectDatabaseUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectDatabaseUserCreate,
		Read:   resourceCloudProjectDatabaseUserRead,
		Delete: resourceCloudProjectDatabaseUserDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectDatabaseUserImportState,
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
			"engine": {
				Type:         schema.TypeString,
				Description:  "Name of the engine of the service",
				ForceNew:     true,
				Required:     true,
				ValidateFunc: helpers.ValidateEnum([]string{"cassandra", "mysql", "kafka", "kafkaConnect"}),
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
	}
}

func resourceCloudProjectDatabaseUserImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	n := 4
	splitId := strings.SplitN(givenId, "/", n)
	if len(splitId) != n {
		return nil, fmt.Errorf("Import Id is not service_name/engine/cluster_id/id formatted")
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

func resourceCloudProjectDatabaseUserCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterId := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/user",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterId),
	)

	params := (&CloudProjectDatabaseUserCreateOpts{}).FromResource(d)
	res := &CloudProjectDatabaseUserResponse{}

	log.Printf("[DEBUG] Will create user: %+v for cluster %s from project %s", params, clusterId, serviceName)
	err := config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		return fmt.Errorf("calling Post %s with params %+v:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for user %s to be READY", res.Id)
	err = waitForCloudProjectDatabaseUserReady(config.OVHClient, serviceName, engine, clusterId, res.Id, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("timeout while waiting user %s to be READY: %w", res.Id, err)
	}
	log.Printf("[DEBUG] user %s is READY", res.Id)

	d.SetId(res.Id)
	d.Set("password", res.Password)

	return resourceCloudProjectDatabaseUserRead(d, meta)
}

func resourceCloudProjectDatabaseUserRead(d *schema.ResourceData, meta interface{}) error {
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
		return helpers.CheckDeleted(d, err, endpoint)
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

func resourceCloudProjectDatabaseUserDelete(d *schema.ResourceData, meta interface{}) error {
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

	log.Printf("[DEBUG] Will delete user %s from cluster %s from project %s", id, clusterId, serviceName)
	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	log.Printf("[DEBUG] Waiting for user %s to be DELETED", id)
	err = waitForCloudProjectDatabaseUserDeleted(config.OVHClient, serviceName, engine, clusterId, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("timeout while waiting user %s to be DELETED: %w", id, err)
	}
	log.Printf("[DEBUG] user %s is DELETED", id)

	d.SetId("")

	return nil
}

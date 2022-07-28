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

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"engine": {
				Type:        schema.TypeString,
				Description: "Name of the engine of the service",
				ForceNew:    true,
				Required:    true,
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
	splitId := strings.SplitN(givenId, "/", 4)
	if len(splitId) != 3 {
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

	err := mustGenericEngineUserEndpoint(engine)
	if err != nil {
		return fmt.Errorf("Calling Post %s :\n\t %q", endpoint, err)
	}

	params := (&CloudProjectDatabaseUserCreateOpts{}).FromResource(d)
	res := &CloudProjectDatabaseUserResponse{}

	log.Printf("[DEBUG] Will create user: %+v for cluster %s from project %s", params, clusterId, serviceName)
	err = config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		return fmt.Errorf("calling Post %s with params %s:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for database %s to be READY", clusterId)
	err = waitForCloudProjectDatabaseReady(config.OVHClient, serviceName, engine, clusterId, 30*time.Second, 5*time.Second)
	if err != nil {
		return fmt.Errorf("timeout while waiting database %s to be READY: %v", clusterId, err)
	}
	log.Printf("[DEBUG] database %s is READY", clusterId)

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

	err := mustGenericEngineUserEndpoint(engine)
	if err != nil {
		return fmt.Errorf("Calling Get %s :\n\t %q", endpoint, err)
	}

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

	err := mustGenericEngineUserEndpoint(engine)
	if err != nil {
		return fmt.Errorf("Calling Delete %s :\n\t %q", endpoint, err)
	}

	log.Printf("[DEBUG] Will delete user %s from cluster %s from project %s", id, clusterId, serviceName)
	err = config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	log.Printf("[DEBUG] Waiting for database %s to be READY", clusterId)
	err = waitForCloudProjectDatabaseReady(config.OVHClient, serviceName, engine, clusterId, 20*time.Second, 1*time.Second)
	if err != nil {
		return fmt.Errorf("timeout while waiting database %s to be READY: %v", clusterId, err)
	}
	log.Printf("[DEBUG] database %s is READY", clusterId)

	d.SetId("")

	return nil
}

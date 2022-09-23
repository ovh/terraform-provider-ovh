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

func resourceCloudProjectDatabaseM3dbUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectDatabaseM3dbUserCreate,
		Read:   resourceCloudProjectDatabaseM3dbUserRead,
		Delete: resourceCloudProjectDatabaseM3dbUserDelete,
		Update: resourceCloudProjectDatabaseM3dbUserUpdate,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectDatabaseM3dbUserImportState,
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
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Id of the database cluster",
				ForceNew:    true,
				Required:    true,
			},
			"group": {
				Type:        schema.TypeString,
				Description: "Group of the user",
				Optional:    true,
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

func resourceCloudProjectDatabaseM3dbUserImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

func resourceCloudProjectDatabaseM3dbUserCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/m3db/%s/user",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
	)
	params := (&CloudProjectDatabaseM3dbUserCreateOpts{}).FromResource(d)
	res := &CloudProjectDatabaseM3dbUserResponse{}

	log.Printf("[DEBUG] Will create user: %+v for cluster %s from project %s", params, clusterId, serviceName)
	err := config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		return fmt.Errorf("calling Post %s with params %+v:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for user %s to be READY", res.Id)
	err = waitForCloudProjectDatabaseUserReady(config.OVHClient, serviceName, "m3db", clusterId, res.Id, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("timeout while waiting user %s to be READY: %w", res.Id, err)
	}
	log.Printf("[DEBUG] user %s is READY", res.Id)

	d.SetId(res.Id)
	d.Set("password", res.Password)

	return resourceCloudProjectDatabaseM3dbUserRead(d, meta)
}

func resourceCloudProjectDatabaseM3dbUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/m3db/%s/user/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)
	res := &CloudProjectDatabaseM3dbUserResponse{}

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

	log.Printf("[DEBUG] Read user %+v", res)
	return nil
}

func resourceCloudProjectDatabaseM3dbUserUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/m3db/%s/user/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)
	params := (&CloudProjectDatabaseM3dbUserUpdateOpts{}).FromResource(d)

	log.Printf("[DEBUG] Will update user: %+v from cluster %s from project %s", params, clusterId, serviceName)
	err := config.OVHClient.Put(endpoint, params, nil)
	if err != nil {
		return fmt.Errorf("calling Put %s with params %+v:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for user %s to be READY", id)
	err = waitForCloudProjectDatabaseUserReady(config.OVHClient, serviceName, "m3db", clusterId, id, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("timeout while waiting user %s to be READY: %w", id, err)
	}
	log.Printf("[DEBUG] user %s is READY", id)

	return resourceCloudProjectDatabaseM3dbUserRead(d, meta)
}

func resourceCloudProjectDatabaseM3dbUserDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/m3db/%s/user/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)

	log.Printf("[DEBUG] Will delete user %s from cluster %s from project %s", id, clusterId, serviceName)
	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	log.Printf("[DEBUG] Waiting for user %s to be DELETED", id)
	err = waitForCloudProjectDatabaseUserDeleted(config.OVHClient, serviceName, "m3db", clusterId, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("timeout while waiting user %s to be DELETED: %w", id, err)
	}
	log.Printf("[DEBUG] user %s is DELETED", id)

	d.SetId("")

	return nil
}

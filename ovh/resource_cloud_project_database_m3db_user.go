package ovh

import (
	"fmt"
	"log"
	"net/url"
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
	}
}

func resourceCloudProjectDatabaseM3dbUserImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return importCloudProjectDatabaseUser(d, meta)
}

func resourceCloudProjectDatabaseM3dbUserCreate(d *schema.ResourceData, meta interface{}) error {
	f := func() interface{} {
		return (&CloudProjectDatabaseM3dbUserCreateOpts{}).FromResource(d)
	}
	return postCloudProjectDatabaseUser(d, meta, "m3db", dataSourceCloudProjectDatabaseM3dbUserRead, resourceCloudProjectDatabaseM3dbUserRead, resourceCloudProjectDatabaseM3dbUserUpdate, f)
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
	f := func() interface{} {
		return (&CloudProjectDatabaseM3dbUserUpdateOpts{}).FromResource(d)
	}
	return updateCloudProjectDatabaseUser(d, meta, "m3db", resourceCloudProjectDatabaseM3dbUserRead, f)
}

func resourceCloudProjectDatabaseM3dbUserDelete(d *schema.ResourceData, meta interface{}) error {
	return deleteCloudProjectDatabaseUser(d, meta, "m3db")
}

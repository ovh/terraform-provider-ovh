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

func resourceHostingPrivateDatabaseDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceHostingPrivateDatabaseDatabaseCreate,
		Read:   resourceHostingPrivateDatabaseDatabaseRead,
		Delete: resourceHostingPrivateDatabaseDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: resourceHostingPrivateDatabaseDatabaseImportState,
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of your private database",
			},
			"database_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of your new database",
			},
		},
	}
}

func resourceHostingPrivateDatabaseDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	databaseName := d.Get("database_name").(string)

	opts := (&HostingPrivateDatabaseDatabaseCreateOpts{}).FromResource(d)
	ds := &HostingPrivateDatabaseDatabase{}

	log.Printf("[DEBUG][Create] HostingPrivateDatabaseDatabase")
	endpoint := fmt.Sprintf("/hosting/privateDatabase/%s/database", url.PathEscape(serviceName))
	err := config.OVHClient.Post(endpoint, opts, &ds)
	if err != nil {
		return fmt.Errorf("failed to create database: %s", err)
	}

	log.Printf("[DEBUG][Create][WaitForArchived] HostingPrivateDatabaseDatabase")
	endpoint = fmt.Sprintf("/hosting/privateDatabase/%s/tasks/%d", url.PathEscape(serviceName), ds.TaskId)
	err = WaitArchivedHostingPrivateDabaseTask(config.OVHClient, endpoint, 2*time.Minute)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", serviceName, databaseName))
	return resourceHostingPrivateDatabaseDatabaseRead(d, meta)
}

func resourceHostingPrivateDatabaseDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	databaseName := d.Get("database_name").(string)

	ds := &HostingPrivateDatabaseDatabase{}

	log.Printf("[DEBUG][Read] HostingPrivateDatabaseDatabase")
	endpoint := fmt.Sprintf("/hosting/privateDatabase/%s/database/%s", url.PathEscape(serviceName), url.PathEscape(databaseName))
	if err := config.OVHClient.Get(endpoint, &ds); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId(fmt.Sprintf("%s/%s", serviceName, databaseName))
	for k, v := range ds.ToMap() {
		d.Set(k, v)
	}

	return nil
}

func resourceHostingPrivateDatabaseDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	databaseName := d.Get("database_name").(string)

	ds := &HostingPrivateDatabaseDatabase{}

	log.Printf("[DEBUG][Delete] HostingPrivateDatabaseDatabase")
	endpoint := fmt.Sprintf("/hosting/privateDatabase/%s/database/%s", url.PathEscape(serviceName), url.PathEscape(databaseName))
	if err := config.OVHClient.Delete(endpoint, ds); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	log.Printf("[DEBUG][Delete][WaitForArchived] HostingPrivateDatabaseDatabase")
	endpoint = fmt.Sprintf("/hosting/privateDatabase/%s/tasks/%d", url.PathEscape(serviceName), ds.TaskId)
	err := WaitArchivedHostingPrivateDabaseTask(config.OVHClient, endpoint, 2*time.Minute)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceHostingPrivateDatabaseDatabaseImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)

	log.Printf("[DEBUG][Import] HostingPrivateDatabaseDatabase givenId: %s", givenId)

	if len(splitId) != 2 {
		return nil, fmt.Errorf("import Id is not SERVICE_NAME/DATABASE_NAME formatted")
	}
	d.SetId(splitId[0])
	d.Set("service_name", splitId[0])
	d.Set("database_name", splitId[1])
	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

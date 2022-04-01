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

func resourceHostingPrivateDatabaseUserGrant() *schema.Resource {
	return &schema.Resource{
		Create: resourceHostingPrivateDatabaseUserGrantCreate,
		Read:   resourceHostingPrivateDatabaseUserGrantRead,
		Delete: resourceHostingPrivateDatabaseUserGrantDelete,
		Importer: &schema.ResourceImporter{
			State: resourceHostingPrivateDatabaseUserGrantImportState,
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of your private database",
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User name used to connect on your databases",
			},
			"database_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Database name where add grant",
			},
			"grant": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Database name where add grant",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateHostingPrivateDatabaseUserGrant(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
		},
	}
}

func resourceHostingPrivateDatabaseUserGrantCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	databaseName := d.Get("database_name").(string)
	userName := d.Get("user_name").(string)
	grant := d.Get("grant").(string)

	opts := (&HostingPrivateDatabaseUserGrantCreateOpts{}).FromResource(d)
	ds := &HostingPrivateDatabaseUserGrant{}

	log.Printf("[DEBUG][Create] HostingPrivateDatabaseUserGrant")
	endpoint := fmt.Sprintf("/hosting/privateDatabase/%s/user/%s/grant", url.PathEscape(serviceName), url.PathEscape(userName))
	err := config.OVHClient.Post(endpoint, opts, &ds)
	if err != nil {
		return fmt.Errorf("failed to create database user grant: %s", err)
	}

	log.Printf("[DEBUG][Create][WaitForArchived] HostingPrivateDatabaseUserGrant")
	endpoint = fmt.Sprintf("/hosting/privateDatabase/%s/tasks/%d", url.PathEscape(serviceName), ds.TaskId)
	err = WaitArchivedHostingPrivateDabaseTask(config.OVHClient, endpoint, 2*time.Minute)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", serviceName, userName, databaseName, grant))
	return resourceHostingPrivateDatabaseUserGrantRead(d, meta)
}

func resourceHostingPrivateDatabaseUserGrantRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	databaseName := d.Get("database_name").(string)
	userName := d.Get("user_name").(string)
	grant := d.Get("grant").(string)

	ds := &HostingPrivateDatabaseUserGrantCreateOpts{}

	log.Printf("[DEBUG][Read] HostingPrivateDatabaseUserGrant")
	endpoint := fmt.Sprintf("/hosting/privateDatabase/%s/user/%s/grant/%s", url.PathEscape(serviceName), url.PathEscape(userName), url.PathEscape(databaseName))
	if err := config.OVHClient.Get(endpoint, &ds); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", serviceName, userName, databaseName, grant))
	for k, v := range ds.ToMap() {
		d.Set(k, v)
	}

	return nil
}

func resourceHostingPrivateDatabaseUserGrantDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	userName := d.Get("user_name").(string)
	databaseName := d.Get("database_name").(string)
	endpoint := fmt.Sprintf("/hosting/privateDatabase/%s/user/%s/grant/%s", url.PathEscape(serviceName), url.PathEscape(userName), url.PathEscape(databaseName))

	ds := &HostingPrivateDatabaseUserGrant{}

	log.Printf("[DEBUG][Delete] HostingPrivateDatabaseUserGrant")
	if err := config.OVHClient.Delete(endpoint, ds); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	log.Printf("[DEBUG][Delete][WaitForArchived] HostingPrivateDatabaseUserGrant")
	endpoint = fmt.Sprintf("/hosting/privateDatabase/%s/tasks/%d", url.PathEscape(serviceName), ds.TaskId)
	err := WaitArchivedHostingPrivateDabaseTask(config.OVHClient, endpoint, 2*time.Minute)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceHostingPrivateDatabaseUserGrantImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 4)

	log.Printf("[DEBUG][Import] HostingPrivateDatabaseUserGrant givenId: %s", givenId)

	if len(splitId) != 4 {
		return nil, fmt.Errorf("import Id is not SERVICE_NAME/USER_NAME/DATABASE_NAME/GRANT formatted")
	}
	d.SetId(splitId[0])
	d.Set("service_name", splitId[0])
	d.Set("user_name", splitId[1])
	d.Set("database_name", splitId[2])
	d.Set("grant", splitId[3])
	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

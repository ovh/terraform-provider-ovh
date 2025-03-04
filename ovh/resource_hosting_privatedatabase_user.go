package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceHostingPrivateDatabaseUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceHostingPrivateDatabaseUserCreate,
		Read:   resourceHostingPrivateDatabaseUserRead,
		Delete: resourceHostingPrivateDatabaseUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceHostingPrivateDatabaseUserImportState,
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of your private database",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "Password for the new user ( alphanumeric and 8 characters minimum )",
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User name used to connect on your databases",
			},
		},
	}
}

func resourceHostingPrivateDatabaseUserCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	userName := d.Get("user_name").(string)

	opts := (&HostingPrivateDatabaseUserCreateOpts{}).FromResource(d)
	ds := &HostingPrivateDatabaseUser{}

	log.Printf("[DEBUG][Create] HostingPrivateDatabaseUser")
	endpoint := fmt.Sprintf("/hosting/privateDatabase/%s/user", url.PathEscape(serviceName))
	err := config.OVHClient.Post(endpoint, opts, &ds)
	if err != nil {
		return fmt.Errorf("failed to create database user: %s", err)
	}

	log.Printf("[DEBUG][Create][WaitForArchived] HostingPrivateDatabaseUser")
	endpoint = fmt.Sprintf("/hosting/privateDatabase/%s/tasks/%d", url.PathEscape(serviceName), ds.TaskId)
	err = WaitArchivedHostingPrivateDabaseTask(config.OVHClient, endpoint, 2*time.Minute)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", serviceName, userName))
	return resourceHostingPrivateDatabaseUserRead(d, meta)
}

func resourceHostingPrivateDatabaseUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	userName := d.Get("user_name").(string)

	ds := &HostingPrivateDatabaseUser{}

	log.Printf("[DEBUG][Read] HostingPrivateDatabaseUser")
	endpoint := fmt.Sprintf("/hosting/privateDatabase/%s/user/%s", url.PathEscape(serviceName), url.PathEscape(userName))
	if err := config.OVHClient.Get(endpoint, &ds); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId(fmt.Sprintf("%s/%s", serviceName, userName))
	for k, v := range ds.ToMap() {
		d.Set(k, v)
	}

	return nil
}

func resourceHostingPrivateDatabaseUserDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	userName := d.Get("user_name").(string)

	ds := &HostingPrivateDatabaseUser{}

	log.Printf("[DEBUG][Delete] HostingPrivateDatabaseUser")
	endpoint := fmt.Sprintf("/hosting/privateDatabase/%s/user/%s", url.PathEscape(serviceName), url.PathEscape(userName))
	if err := config.OVHClient.Delete(endpoint, ds); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	log.Printf("[DEBUG][Delete][WaitForArchived] HostingPrivateDatabaseUser")
	endpoint = fmt.Sprintf("/hosting/privateDatabase/%s/tasks/%d", url.PathEscape(serviceName), ds.TaskId)
	err := WaitArchivedHostingPrivateDabaseTask(config.OVHClient, endpoint, 2*time.Minute)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceHostingPrivateDatabaseUserImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)

	log.Printf("[DEBUG][Import] HostingPrivateDatabaseUser givenId: %s", givenId)

	if len(splitId) != 2 {
		return nil, fmt.Errorf("import Id is not SERVICE_NAME/USER_NAME formatted")
	}
	d.SetId(splitId[0])
	d.Set("service_name", splitId[0])
	d.Set("user_name", splitId[1])
	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

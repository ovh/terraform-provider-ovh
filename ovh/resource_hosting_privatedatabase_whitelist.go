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

func resourceHostingPrivateDatabaseWhitelist() *schema.Resource {
	return &schema.Resource{
		Create: resourceHostingPrivateDatabaseWhitelistCreate,
		Update: resourceHostingPrivateDatabaseWhitelistUpdate,
		Read:   resourceHostingPrivateDatabaseWhitelistRead,
		Delete: resourceHostingPrivateDatabaseWhitelistDelete,
		Importer: &schema.ResourceImporter{
			State: resourceHostingPrivateDatabaseWhitelistImportState,
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your private database",
			},
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The whitelisted IP in your instance",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom name for your Whitelisted IP",
			},
			"service": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Authorize this IP to access service port",
			},
			"sftp": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Authorize this IP to access SFTP port",
			},
		},
	}
}

func resourceHostingPrivateDatabaseWhitelistCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	ip := HostingPrivateDatabaseWhitelistefaultNetmask(d.Get("ip").(string))

	opts := (&HostingPrivateDatabaseWhitelistCreateOpts{}).FromResource(d)
	ds := &HostingPrivateDatabaseWhitelist{}

	log.Printf("[DEBUG][Create] HostingPrivateDatabaseWhitelist")
	endpoint := fmt.Sprintf("/hosting/privateDatabase/%s/whitelist", url.PathEscape(serviceName))
	err := config.OVHClient.Post(endpoint, opts, &ds)
	if err != nil {
		return fmt.Errorf("failed to create HostingPrivateDatabaseWhitelist: %s", err)
	}

	log.Printf("[DEBUG][Create][WaitForArchived] HostingPrivateDatabaseWhitelist")
	endpoint = fmt.Sprintf("/hosting/privateDatabase/%s/tasks/%d", url.PathEscape(serviceName), ds.TaskId)
	err = WaitArchivedHostingPrivateDabaseTask(config.OVHClient, endpoint, 2*time.Minute)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", serviceName, ip))
	return resourceHostingPrivateDatabaseWhitelistRead(d, meta)
}

func resourceHostingPrivateDatabaseWhitelistUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	ip := HostingPrivateDatabaseWhitelistefaultNetmask(d.Get("ip").(string))

	opts := (&HostingPrivateDatabaseWhitelistUpdateOpts{}).FromResource(d)

	log.Printf("[DEBUG][Update] HostingPrivateDatabaseWhitelist")
	endpoint := fmt.Sprintf("/hosting/privateDatabase/%s/whitelist/%s", url.PathEscape(serviceName), url.PathEscape(ip))
	err := config.OVHClient.Put(endpoint, opts, nil)
	if err != nil {
		return fmt.Errorf("failed to update whitelist: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", serviceName, ip))
	return resourceHostingPrivateDatabaseWhitelistRead(d, meta)
}

func resourceHostingPrivateDatabaseWhitelistRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	ip := HostingPrivateDatabaseWhitelistefaultNetmask(d.Get("ip").(string))

	ds := &HostingPrivateDatabaseWhitelist{}

	log.Printf("[DEBUG][Read] HostingPrivateDatabaseWhitelist")
	endpoint := fmt.Sprintf("/hosting/privateDatabase/%s/whitelist/%s", url.PathEscape(serviceName), url.PathEscape(ip))
	if err := config.OVHClient.Get(endpoint, &ds); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId(fmt.Sprintf("%s/%s", serviceName, ip))
	for k, v := range ds.ToMap() {
		d.Set(k, v)
	}

	return nil
}

func resourceHostingPrivateDatabaseWhitelistDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	ip := HostingPrivateDatabaseWhitelistefaultNetmask(d.Get("ip").(string))

	ds := &HostingPrivateDatabaseWhitelist{}

	log.Printf("[DEBUG][Delete] HostingPrivateDatabaseWhitelist")
	endpoint := fmt.Sprintf("/hosting/privateDatabase/%s/whitelist/%s", url.PathEscape(serviceName), url.PathEscape(ip))
	if err := config.OVHClient.Delete(endpoint, ds); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	log.Printf("[DEBUG][Delete][WaitForArchived] HostingPrivateDatabaseWhitelist")
	endpoint = fmt.Sprintf("/hosting/privateDatabase/%s/tasks/%d", url.PathEscape(serviceName), ds.TaskId)
	err := WaitArchivedHostingPrivateDabaseTask(config.OVHClient, endpoint, 2*time.Minute)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceHostingPrivateDatabaseWhitelistImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)

	log.Printf("[DEBUG][Import] HostingPrivateDatabaseWhitelist givenId: %s", givenId)

	if len(splitId) != 2 {
		return nil, fmt.Errorf("import Id is not SERVICE_NAME/IP formatted")
	}
	d.SetId(splitId[0])
	d.Set("service_name", splitId[0])
	d.Set("ip", splitId[1])
	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

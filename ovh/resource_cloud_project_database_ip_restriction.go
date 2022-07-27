package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func resourceCloudProjectDatabaseIpRestriction() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectDatabaseIpRestrictionCreate,
		Read:   resourceCloudProjectDatabaseIpRestrictionRead,
		Delete: resourceCloudProjectDatabaseIpRestrictionDelete,
		Update: resourceCloudProjectDatabaseIpRestrictionUpdate,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectDatabaseIpRestrictionImportState,
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
			"ip": {
				Type:        schema.TypeString,
				Description: "Authorized IP",
				ForceNew:    true,
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the IP restriction",
				Optional:    true,
			},

			//Computed
			"status": {
				Type:        schema.TypeString,
				Description: "Current status of the IP restriction",
				Computed:    true,
			},
		},
	}
}

func resourceCloudProjectDatabaseIpRestrictionImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	n := 4
	splitId := strings.SplitN(givenId, "/", n)
	if len(splitId) != n {
		return nil, fmt.Errorf("Import Id is not service_name/engine/cluster_id/ip formatted")
	}
	serviceName := splitId[0]
	engine := splitId[1]
	clusterId := splitId[2]
	ip := splitId[3]
	d.SetId(hashcode.Strings([]string{ip}))
	d.Set("engine", engine)
	d.Set("service_name", serviceName)
	d.Set("cluster_id", clusterId)
	d.Set("ip", ip)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectDatabaseIpRestrictionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterId := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/ipRestriction",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterId),
	)
	params := (&CloudProjectDatabaseIpRestrictionCreateOpts{}).FromResource(d)
	res := &CloudProjectDatabaseIpRestrictionResponse{}

	log.Printf("[DEBUG] Will create IP restriction: %+v for cluster %s from project %s", params, clusterId, serviceName)
	err := config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		return fmt.Errorf("calling Post %s with params %s:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for database %s to be READY", clusterId)
	err = waitForCloudProjectDatabaseReady(config.OVHClient, serviceName, engine, clusterId, 20*time.Second, 1*time.Second)
	if err != nil {
		return fmt.Errorf("timeout while waiting database %s to be READY: %v", clusterId, err)
	}
	log.Printf("[DEBUG] database %s is READY", clusterId)

	d.SetId(hashcode.Strings([]string{res.Ip}))

	return resourceCloudProjectDatabaseIpRestrictionRead(d, meta)
}

func resourceCloudProjectDatabaseIpRestrictionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterId := d.Get("cluster_id").(string)
	ip := d.Get("ip").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/ipRestriction/%s",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterId),
		url.PathEscape(ip),
	)
	res := &CloudProjectDatabaseIpRestrictionResponse{}

	log.Printf("[DEBUG] Will read IP restriction %s from cluster %s from project %s", ip, clusterId, serviceName)
	if err := config.OVHClient.Get(endpoint, res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId(hashcode.Strings([]string{res.Ip}))
	for k, v := range res.ToMap() {
		d.Set(k, v)
	}

	log.Printf("[DEBUG] Read IP restriction %+v", res)
	return nil
}

func resourceCloudProjectDatabaseIpRestrictionUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterId := d.Get("cluster_id").(string)
	ip := d.Get("ip").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/ipRestriction/%s",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterId),
		url.PathEscape(ip),
	)
	params := (&CloudProjectDatabaseIpRestrictionUpdateOpts{}).FromResource(d)

	log.Printf("[DEBUG] Will update IP restriction: %+v from cluster %s from project %s", params, clusterId, serviceName)
	err := config.OVHClient.Put(endpoint, params, nil)
	if err != nil {
		return fmt.Errorf("calling Put %s with params %v:\n\t %q", endpoint, params, err)
	}

	return resourceCloudProjectDatabaseIpRestrictionRead(d, meta)
}

func resourceCloudProjectDatabaseIpRestrictionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	clusterId := d.Get("cluster_id").(string)
	ip := d.Get("ip").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/ipRestriction/%s",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterId),
		url.PathEscape(ip),
	)

	log.Printf("[DEBUG] Will delete IP restriction %s from cluster %s from project %s", ip, clusterId, serviceName)
	err := config.OVHClient.Delete(endpoint, nil)
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

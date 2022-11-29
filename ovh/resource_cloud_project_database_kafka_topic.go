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

func resourceCloudProjectDatabaseKafkaTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectDatabaseKafkaTopicCreate,
		Read:   resourceCloudProjectDatabaseKafkaTopicRead,
		Delete: resourceCloudProjectDatabaseKafkaTopicDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectDatabaseKafkaTopicImportState,
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
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the topic",
				ForceNew:    true,
				Required:    true,
			},

			//Optional/Computed
			"min_insync_replicas": {
				Type:         schema.TypeInt,
				Description:  "Minimum insync replica accepted for this topic",
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateCloudProjectDatabaseKafkaTopicMinInsyncReplicasFunc,
			},
			"partitions": {
				Type:         schema.TypeInt,
				Description:  "Number of partitions for this topic",
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateCloudProjectDatabaseKafkaTopicPartitionsFunc,
			},
			"replication": {
				Type:         schema.TypeInt,
				Description:  "Number of replication for this topic",
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateCloudProjectDatabaseKafkaTopicReplicationFunc,
			},
			"retention_bytes": {
				Type:        schema.TypeInt,
				Description: "Number of bytes for the retention of the data for this topic",
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
			"retention_hours": {
				Type:         schema.TypeInt,
				Description:  "Number of hours for the retention of the data for this topic",
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateCloudProjectDatabaseKafkaTopicRetentionHoursFunc,
			},
		},
	}
}

func resourceCloudProjectDatabaseKafkaTopicImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

func resourceCloudProjectDatabaseKafkaTopicCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/topic",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
	)
	params := (&CloudProjectDatabaseKafkaTopicCreateOpts{}).FromResource(d)
	res := &CloudProjectDatabaseKafkaTopicResponse{}

	log.Printf("[DEBUG] Will create topic: %+v for cluster %s from project %s", params, clusterId, serviceName)
	err := config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Waiting for topic %s to be READY", res.Id)
	err = waitForCloudProjectDatabaseKafkaTopicReady(config.OVHClient, serviceName, clusterId, res.Id, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("timeout while waiting topic %s to be READY: %w", res.Id, err)
	}
	log.Printf("[DEBUG] topic %s is READY", res.Id)

	d.SetId(res.Id)

	err = resourceCloudProjectDatabaseKafkaTopicRead(d, meta)
	if err != nil {
		return err
	}
	return nil
}

func resourceCloudProjectDatabaseKafkaTopicRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/topic/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)
	res := &CloudProjectDatabaseKafkaTopicResponse{}

	log.Printf("[DEBUG] Will read topic %s from cluster %s from project %s", id, clusterId, serviceName)
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

	log.Printf("[DEBUG] Read topic %+v", res)
	return nil
}

func resourceCloudProjectDatabaseKafkaTopicDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/topic/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)

	log.Printf("[DEBUG] Will delete topic  %s from cluster %s from project %s", id, clusterId, serviceName)
	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		err = helpers.CheckDeleted(d, err, endpoint)
		if err != nil {
			return err
		}
		return nil
	}

	log.Printf("[DEBUG] Waiting for topic %s to be DELETED", id)
	err = waitForCloudProjectDatabaseKafkaTopicDeleted(config.OVHClient, serviceName, clusterId, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("timeout while waiting topic %s to be DELETED: %w", id, err)
	}
	log.Printf("[DEBUG] topic %s is DELETED", id)

	d.SetId("")

	return nil
}

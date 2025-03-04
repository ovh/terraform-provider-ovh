package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

const (
	OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_ENDPOINT = "/cloud/project/%s/region/%s/workflow/backup"
	OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_REGION   = "region_name"
	OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_SERVICE  = "service_name"
)

func resourceCloudProjectWorkflowBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectWorkflowBackupCreate,
		Read:   resourceCloudProjectWorkflowBackupRead,
		Delete: resourceCloudProjectWorkflowBackupDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: resourceCloudProjectWorkflowBackupSchema(),
	}
}

// No endpoint to update the workflow hence all attributes are ForceNew : true
func resourceCloudProjectWorkflowBackupSchema() map[string]*schema.Schema {
	schema := map[string]*schema.Schema{
		OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_SERVICE: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
		},
		OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_REGION: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Region name.",
		},
		"cron": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"instance_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"max_execution_count": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: false,
			ForceNew: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"rotation": {
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
		},

		"backup_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		// computed
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
	return schema
}

func resourceCloudProjectWorkflowBackupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get(OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_SERVICE).(string)
	regionName := d.Get(OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_REGION).(string)
	endpoint := fmt.Sprintf(OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_ENDPOINT, serviceName, regionName)
	params := (&CloudProjectWorkflowBackupCreateOpts{}).FromResource(d)
	res := &CloudProjectWorkflowBackupResponse{}

	log.Printf("[DEBUG] will create workflow on %s", endpoint)
	err := config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		return fmt.Errorf("calling Post %s with param %+v:\n\t %w", endpoint, params, err)
	}
	d.SetId(res.Id)
	return resourceCloudProjectWorkflowBackupRead(d, meta)
}

func resourceCloudProjectWorkflowBackupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get(OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_SERVICE).(string)
	regionName := d.Get(OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_REGION).(string)
	endpoint := fmt.Sprintf(OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_ENDPOINT, serviceName, regionName) + "/" + d.Id()
	res := &CloudProjectWorkflowBackupResponse{}

	log.Printf("[DEBUG] Will read workflow backup on %s", endpoint)
	if err := config.OVHClient.Get(endpoint, res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}
	log.Printf("[DEBUG] Response %s", res.ToMap())

	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, fmt.Sprint(v))
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	log.Printf("[DEBUG] read  %+v", res)
	return nil
}

func resourceCloudProjectWorkflowBackupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get(OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_SERVICE).(string)
	regionName := d.Get(OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_REGION).(string)
	endpoint := fmt.Sprintf(OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_ENDPOINT, serviceName, regionName) + "/" + d.Id()
	log.Printf("[DEBUG] Will delete workflow backup on %s", endpoint)
	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}
	return nil
}

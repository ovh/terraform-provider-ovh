package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceDedicatedServerInstallTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceDedicatedServerInstallTaskCreate,
		Update: resourceDedicatedServerInstallTaskUpdate,
		Read:   resourceDedicatedServerInstallTaskRead,
		Delete: resourceDedicatedServerInstallTaskDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDedicatedServerInstallTaskImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of your dedicated server.",
			},
			"partition_scheme_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Partition scheme name.",
			},
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Template name",
			},
			"bootid_on_destroy": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "If set, reboot the server on the specified boot id during destroy phase",
			},
			"details": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"custom_hostname": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Set up the server using the provided hostname instead of the default hostname",
						},
						"disk_group_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "",
						},
						"install_sql_server": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: "",
						},
						"language": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "language",
						},
						"no_raid": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: "",
						},
						"post_installation_script_link": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Indicate the URL where your postinstall customisation script is located",
						},
						"post_installation_script_return": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "indicate the string returned by your postinstall customisation script on successful execution. Advice: your script should return a unique validation string in case of succes. A good example is 'loh1Xee7eo OK OK OK UGh8Ang1Gu'",
						},
						"soft_raid_devices": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "",
						},
						"ssh_key_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Name of the ssh key that should be installed. Password login will be disabled",
						},
						"use_spla": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: "",
						},
					},
				},
			},
			"user_metadata": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 128,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Description: "The key for the user_metadata",
							Optional:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "The value for the user_metadata",
							Optional:    true,
						},
					},
				},
			},

			//Computed
			"comment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Details of this task",
			},
			"done_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Completion date",
			},
			"function": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Function name",
			},
			"last_update": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update",
			},
			"start_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task Creation date",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task status",
			},
		},
	}
}

func resourceDedicatedServerInstallTaskImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	//
	// After creating an install task, there is no way to get the name of the server (service_name)
	// nor the name of template (template_name) used during the creation of the task.
	// This is why it is required to provide them when importing.
	//

	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 3)
	if len(splitId) != 3 {
		return nil, fmt.Errorf("Import Id is not service_name/template_name/task_id formatted")
	}
	serviceName := splitId[0]
	templateName := splitId[1]
	taskId := splitId[2]
	d.SetId(taskId)
	d.Set("service_name", serviceName)
	d.Set("template_name", templateName)
	err := dedicatedServerInstallTaskRead(d, meta)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, err
}

func resourceDedicatedServerInstallTaskCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf(
		"/dedicated/server/%s/install/start",
		url.PathEscape(serviceName),
	)
	opts := (&DedicatedServerInstallTaskCreateOpts{}).FromResource(d)
	task := &DedicatedServerTask{}

	if err := config.OVHClient.Post(endpoint, opts, task); err != nil {
		// POST on install tasks can fail randomly so in order to avoid issues, let's allow
		// a retry via waitForDedicatedServerTask
		log.Printf("[WARN] Ignored error when calling POST %s: %v", endpoint, err)
	}

	if err := waitForDedicatedServerTask(serviceName, task, config.OVHClient); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", task.Id))

	return dedicatedServerInstallTaskRead(d, meta)
}

func dedicatedServerInstallTaskRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf(
			"Could not parse install task id %s,%s:\n\t %q",
			serviceName,
			d.Id(),
			err,
		)
	}

	task, err := getDedicatedServerTask(serviceName, id, config.OVHClient)
	if err != nil {
		return helpers.CheckDeleted(d, err, fmt.Sprintf(
			"dedicated server task %s/%s",
			serviceName,
			d.Id(),
		))
	}

	d.Set("function", task.Function)
	d.Set("comment", task.Comment)
	d.Set("status", task.Status)
	d.Set("last_update", task.LastUpdate.Format(time.RFC3339))
	d.Set("done_date", task.DoneDate.Format(time.RFC3339))
	d.Set("start_date", task.StartDate.Format(time.RFC3339))

	return nil
}

func resourceDedicatedServerInstallTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	// nothing to do on update
	return resourceDedicatedServerInstallTaskRead(d, meta)
}

func resourceDedicatedServerInstallTaskRead(d *schema.ResourceData, meta interface{}) error {
	// Nothing to do on READ
	//
	// IMPORTANT: This resource doesn't represent a real resource
	// but instead a task on a dedicated server. OVH may clean its tasks database after a while
	// so that the API may return a 404 on a task id. If we hit a 404 on a READ, then
	// terraform will understand that it has to recreate the resource, and consequently
	// will trigger new install task on the dedicated server.
	// This is something we must avoid!
	//
	return nil
}

func resourceDedicatedServerInstallTaskDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	bootId := helpers.GetNilIntPointerFromData(d, "bootid_on_destroy")

	if bootId != nil {
		serviceName := d.Get("service_name").(string)

		// before reboot, update bootId accordingly
		bootIdEndpoint := fmt.Sprintf("/dedicated/server/%s", url.PathEscape(serviceName))
		bootIdReqBody := make(map[string]int)
		bootIdReqBody["bootId"] = *bootId
		if err := config.OVHClient.Put(bootIdEndpoint, bootIdReqBody, nil); err != nil {
			return fmt.Errorf("Error calling PUT %s:\n\t %q", bootIdEndpoint, err)
		}

		// reboot
		endpoint := fmt.Sprintf(
			"/dedicated/server/%s/reboot",
			url.PathEscape(serviceName),
		)

		task := &DedicatedServerTask{}
		if err := config.OVHClient.Post(endpoint, nil, task); err != nil {
			// POST on install tasks can fail randomly so in order to avoid issues, let's allow
			// a retry via waitForDedicatedServerTask
			log.Printf("[WARN] Ignored error when calling POST %s: %v", endpoint, err)
		}

		if err := waitForDedicatedServerTask(serviceName, task, config.OVHClient); err != nil {
			return err
		}
	}

	// we cant delete the task through the API, just forget about its Id
	d.SetId("")
	return nil
}

package ovh

import (
	"fmt"
	"net/url"
	"strconv"
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
						"change_log": {
							Type:        schema.TypeString,
							Optional:    true,
							Deprecated:  "field is not used anymore",
							Description: "Template change log details",
						},
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
						"install_rtm": {
							Type:        schema.TypeBool,
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
						"reset_hw_raid": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: "",
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
						"use_distrib_kernel": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: "Use the distribution's native kernel instead of the recommended OVH Kernel",
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
		return fmt.Errorf("Error calling POST %s:\n\t %q", endpoint, err)
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
		endpoint := fmt.Sprintf(
			"/dedicated/server/%s/reboot",
			url.PathEscape(serviceName),
		)

		task := &DedicatedServerTask{}
		if err := config.OVHClient.Post(endpoint, nil, task); err != nil {
			return fmt.Errorf("Error calling POST %s:\n\t %q", endpoint, err)
		}

		if err := waitForDedicatedServerTask(serviceName, task, config.OVHClient); err != nil {
			return err
		}
	}

	// we cant delete the task through the API, just forget about its Id
	d.SetId("")
	return nil
}

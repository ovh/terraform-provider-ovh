package ovh

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceDedicatedServerBringYourOwnImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceDedicatedServerBringYourOwnImageCreate,
		Update: resourceDedicatedServerBringYourOwnImageUpdate,
		Read:   resourceDedicatedServerBringYourOwnImageTaskRead,
		Delete: resourceDedicatedServerBringYourOwnImageTaskDelete,

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
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Image URL",
				ValidateDiagFunc: func(value interface{}, attributePath cty.Path) diag.Diagnostics {
					if _, err := url.Parse(value.(string)); err != nil {
						return diag.Diagnostics{
							{
								Severity:      diag.Error,
								Summary:       "Image URL is invalid",
								Detail:        fmt.Sprintf("Image URL '%v' is invalid: %v", value, err),
								AttributePath: attributePath,
							},
						}
					}
					return diag.Diagnostics{}
				},
			},
			"check_sum": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Image checksum",
			},
			"check_sum_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Checksum type",
			},
			"config_drive": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Set: func(i interface{}) int {
					out := fmt.Sprintf("%#v", i)
					hash := int(schema.HashString(out))
					return hash
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"hostname": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"ssh_key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"user_data": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"user_metadata": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							/*Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type: schema.TypeString,
									},
									"value": {
										Type: schema.TypeString,
									},
								},
							},*/
						},
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Image description",
				Optional:    true,
			},
			"disk_group_id": {
				Type:        schema.TypeFloat,
				Description: "Disk group id to process install on (only available for some templates)",
				Optional:    true,
			},
			"http_headers": {
				Description: "HTTP Headers",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Image type",
				Required:    true,
				ValidateDiagFunc: func(value interface{}, attributePath cty.Path) diag.Diagnostics {
					if value != "qcow2" && value != "raw" {
						return diag.Diagnostics{
							{
								Severity:      diag.Error,
								Summary:       "Image type is not set or is invalid",
								Detail:        fmt.Sprintf("Image type is currently set as '%v'. It should be either 'qcow2' or 'raw'.", value),
								AttributePath: attributePath,
							},
						}
					}
					return diag.Diagnostics{}
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
				Description: "Bring Your Own Image status",
			},
		},
	}
}

func resourceDedicatedServerBringYourOwnImageCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf(
		"/dedicated/server/%s/bringYourOwnImage",
		url.PathEscape(serviceName),
	)
	opts, err := (&DedicatedServerBringYourOwnImageCreateOpts{}).FromResource(d)
	if err != nil {
		return err
	}
	task := &DedicatedServerTask{}

	if err := config.OVHClient.Post(endpoint, opts, &task); err != nil {
		return fmt.Errorf("error calling POST %s:\n\t %q", endpoint, err)
	}

	if task == nil || task.Id == 0 {
		task, err = getLatestDedicatedServerTasks(serviceName, "reinstallServer", config.OVHClient)
		if err != nil {
			return err
		}
	}

	if err := waitForDedicatedServerTask(serviceName, task, config.OVHClient); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", task.Id))

	return dedicatedServerBringYourOwnImageTaskCreate(d, meta)
}

func dedicatedServerBringYourOwnImageTaskCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf(
			"could not parse BringYourOwnImage task id %s,%s:\n\t %q",
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

func resourceDedicatedServerBringYourOwnImageUpdate(d *schema.ResourceData, meta interface{}) error {
	// nothing to do on update
	return nil
}

func resourceDedicatedServerBringYourOwnImageTaskRead(d *schema.ResourceData, meta interface{}) error {
	// Nothing to do on READ
	//
	// IMPORTANT: This resource doesn't represent a real resource
	// but instead a task on a dedicated server. OVH may clean its tasks database after a while
	// so that the API may return a 404 on a task id. If we hit a 404 on a READ, then
	// terraform will understand that it has to recreate the resource, and consequently
	// will trigger new reinstallServer task on the dedicated server.
	// This is something we must avoid!
	//
	return nil
}

func resourceDedicatedServerBringYourOwnImageTaskDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf(
		"/dedicated/server/%s/bringYourOwnImage",
		url.PathEscape(serviceName),
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return fmt.Errorf("error calling DELETE %s:\n\t %q", endpoint, err)
	}

	// we cant delete the task through the API, just forget about its Id
	d.SetId("")
	return nil
}

func getLatestDedicatedServerTasks(serviceName, function string, c *ovh.Client) (*DedicatedServerTask, error) {
	tasksEndpoint := fmt.Sprintf("/dedicated/server/%s/task?function=%s", url.PathEscape(serviceName), function)

	taskIds := []int{}
	if err := c.Get(tasksEndpoint, &taskIds); err != nil {
		return nil, err
	}

	if len(taskIds) == 0 {
		return nil, nil
	}

	taskEndpoint := fmt.Sprintf("/dedicated/server/%s/task/%d", url.PathEscape(serviceName), taskIds[0])

	task := &DedicatedServerTask{}
	if err := c.Get(taskEndpoint, &task); err != nil {
		return nil, err
	}

	return task, nil
}

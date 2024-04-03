package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_url": {
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
						"http_headers_0_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Required:    false,
							Description: "httpHeaders0Key",
						},
						"http_headers_0_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Required:    false,
							Description: "httpHeaders0Value",
						},
						"http_headers_1_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Required:    false,
							Description: "httpHeaders1Key",
						},
						"http_headers_1_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Required:    false,
							Description: "httpHeaders1Value",
						},
						"http_headers_2_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Required:    false,
							Description: "httpHeaders2Key",
						},
						"http_headers_2_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Required:    false,
							Description: "httpHeaders2Value",
						},
						"http_headers_3_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Required:    false,
							Description: "httpHeaders3Key",
						},
						"http_headers_3_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Required:    false,
							Description: "httpHeaders3Value",
						},
						"http_headers_4_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Required:    false,
							Description: "httpHeaders4Key",
						},
						"http_headers_4_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Required:    false,
							Description: "httpHeaders4Value",
						},
						"http_headers_5_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Required:    false,
							Description: "httpHeaders5Key",
						},
						"http_headers_5_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Required:    false,
							Description: "httpHeaders5Value",
						},
						"image_checksum": {
							Type:        schema.TypeString,
							Optional:    true,
							Required:    false,
							Description: "Image checksum",
						},
						"image_checksum_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Required:    false,
							Description: "Checksum type",
							ValidateDiagFunc: func(value interface{}, attributePath cty.Path) diag.Diagnostics {
								if value != "md5" &&
									value != "sha1" &&
									value != "sha256" &&
									value != "sha512" {
									return diag.Diagnostics{
										{
											Severity:      diag.Error,
											Summary:       "Checksum type is not set or is invalid",
											Detail:        fmt.Sprintf("Checksum type is currently set as '%v'. It should be one of 'md5', 'sha1', 'sha256' or 'sha512'.", value),
											AttributePath: attributePath,
										},
									}
								}
								return diag.Diagnostics{}
							},
						},
						"config_drive_user_data": {
							Type:        schema.TypeString,
							Optional:    true,
							Required:    false,
							Description: "configDriveUserData",
						},
						"image_type": {
							Type:        schema.TypeString,
							Description: "Image type",
							Optional:    true,
							Required:    false,
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
						"language": {
							Type:        schema.TypeString,
							Optional:    true,
							Required:    false,
							Description: "language",
						},
						"use_spla": {
							Type:        schema.TypeBool,
							Optional:    true,
							Required:    false,
							Description: "set to true to use your own licence",
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

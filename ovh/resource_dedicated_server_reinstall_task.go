package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceDedicatedServerReinstallTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceDedicatedServerReinstallTaskCreate,
		Update: resourceDedicatedServerReinstallTaskUpdate,
		Read:   resourceDedicatedServerReinstallTaskRead,
		Delete: resourceDedicatedServerReinstallTaskDelete,

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
			"os": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Operating System name",
			},
			"bootid_on_destroy": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "If set, reboot the server on the specified boot id during destroy phase",
			},
			"customizations": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "OS reinstallation customizations",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_drive_metadata": {
							Type:        schema.TypeMap,
							Optional:    true,
							ForceNew:    true,
							Description: "Config Drive MetaData",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"config_drive_user_data": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Config Drive UserData",
						},
						"efi_bootloader_path": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "EFI bootloader path",
						},
						"hostname": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Custom hostname",
						},
						"http_headers": {
							Type:        schema.TypeMap,
							Optional:    true,
							ForceNew:    true,
							Description: "Image HTTP Headers",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"image_check_sum": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Image checksum",
						},
						"image_check_sum_type": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Checksum type",
						},
						"image_type": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Image Type",
						},
						"image_url": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Image URL",
						},
						"language": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Display Language",
						},
						"post_installation_script": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Post-Installation Script",
						},
						"post_installation_script_extension": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Post-Installation Script File Extension",
						},
						"ssh_key": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "SSH Public Key",
						},
					},
				},
			},
			"properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Attribute 'properties' is deprecated and has no effect",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"storage": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Storage configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_group_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Disk group id (default is 0, meaning automatic)",
						},
						"hardware_raid": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "Hardware Raid configurations (if not specified, all disks of the chosen disk group id will be configured in JBOD mode)",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"arrays": {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    true,
										Description: "Number of arrays (default is 1)",
									},
									"disks": {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    true,
										Description: "Total number of disks in the disk group involved in the hardware raid configuration (all disks of the disk group by default)",
									},
									"raid_level": {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    true,
										Description: "Hardware raid type (default is 1)",
									},
									"spares": {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    true,
										Description: "Number of disks in the disk group involved in the spare (default is 0)",
									},
								},
							},
						},
						"partitioning": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "Partitioning configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scheme_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Partitioning scheme name",
									},
									"layout": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Custom partitioning layout (default is the default layout of the operating system's default partitioning scheme)",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"extras": {
													Type:        schema.TypeList,
													Optional:    true,
													ForceNew:    true,
													Description: "Partition extras parameters",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"lv": {
																Type:        schema.TypeList,
																Optional:    true,
																ForceNew:    true,
																Description: "LVM-specific parameters",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			ForceNew:    true,
																			Description: "Logical volume name",
																		},
																	},
																},
															},
															"zp": {
																Type:        schema.TypeList,
																Optional:    true,
																ForceNew:    true,
																Description: "ZFS-specific parameters",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			ForceNew:    true,
																			Description: "zpool name (generated automatically if not specified)",
																		},
																	},
																},
															},
														},
													},
												},
												"file_system": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "File system type",
												},
												"mount_point": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Mount point",
												},
												"raid_level": {
													Type:        schema.TypeInt,
													Default:     1,
													Optional:    true,
													ForceNew:    true,
													Description: "Software raid type (default is 1)",
												},
												"size": {
													Type:        schema.TypeInt,
													Optional:    true,
													ForceNew:    true,
													Description: "Partition size in MiB (default value is 0)",
												},
											},
										},
									},
									"disks": {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    true,
										Description: "Total number of disks in the disk group involved in the partitioning configuration (all disks of the disk group by default)",
									},
								},
							},
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

func resourceDedicatedServerReinstallTaskCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf(
		"/dedicated/server/%s/reinstall",
		url.PathEscape(serviceName),
	)
	opts := (&DedicatedServerReinstallTaskCreateOpts{}).FromResource(d)
	task := DedicatedServerTask{}

	if err := config.OVHClient.Post(endpoint, opts, &task); err != nil {
		// If task was not created because of an error, return it immediately.
		if task.Id == 0 {
			return fmt.Errorf("failed to create reinstall task: %w", err)
		}

		// POST on reinstall tasks can fail randomly so in order to avoid issues, let's allow
		// a retry via waitForDedicatedServerTask
		log.Printf("[WARN] Ignored error when calling POST %s: %v", endpoint, err)
	}

	if err := waitForDedicatedServerTask(serviceName, &task, config.OVHClient); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", task.Id))

	return dedicatedServerReinstallTaskRead(d, meta)
}

func dedicatedServerReinstallTaskRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf(
			"Could not parse reinstall task id %s,%s:\n\t %q",
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

func resourceDedicatedServerReinstallTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	// nothing to do on update
	return resourceDedicatedServerReinstallTaskRead(d, meta)
}

func resourceDedicatedServerReinstallTaskRead(d *schema.ResourceData, meta interface{}) error {
	// Nothing to do on READ
	//
	// IMPORTANT: This resource doesn't represent a real resource
	// but instead a task on a dedicated server. OVH may clean its tasks database after a while
	// so that the API may return a 404 on a task id. If we hit a 404 on a READ, then
	// terraform will understand that it has to recreate the resource, and consequently
	// will trigger new reinstall task on the dedicated server.
	// This is something we must avoid!
	//
	return nil
}

func resourceDedicatedServerReinstallTaskDelete(d *schema.ResourceData, meta interface{}) error {
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
			// POST on reinstall tasks can fail randomly so in order to avoid issues, let's allow
			// a retry via waitForDedicatedServerTask
			// If task was not created because of an error, return it immediately.
			if task.Id == 0 {
				return fmt.Errorf("failed to create reboot task: %w", err)
			}
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

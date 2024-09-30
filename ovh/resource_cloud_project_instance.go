package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceCloudProjectInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudProjectInstanceCreate,
		ReadContext:   resourceCloudProjectInstanceRead,
		DeleteContext: resourceCloudProjectInstanceDelete,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
				Description: "Service name of the resource representing the id of the cloud project.",
				ForceNew:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "Instance region",
				Required:    true,
				ForceNew:    true,
			},
			"auto_backup": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "Create an autobackup workflow after instance start up",
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cron": {
							Type:        schema.TypeString,
							Description: "Unix cron pattern",
							ForceNew:    true,
							Required:    true,
						},
						"rotation": {
							Type:        schema.TypeInt,
							Description: "Number of backup to keep",
							ForceNew:    true,
							Required:    true,
						},
					},
				},
			},
			"billing_period": {
				Type:         schema.TypeString,
				Description:  "Billing period - hourly | monthly ",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: helpers.ValidateEnum([]string{"monthly", "hourly"}),
			},
			"boot_from": {
				Type:        schema.TypeSet,
				Required:    true,
				MaxItems:    1,
				Description: "Boot the instance from an image or a volume",
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_id": {
							Type:        schema.TypeString,
							Description: "Instance image id",
							Optional:    true,
						},
						"volume_id": {
							Type:        schema.TypeString,
							Description: "Instance volume id",
							Optional:    true,
						},
					},
				},
			},
			"bulk": {
				Type:        schema.TypeInt,
				Description: "Create multiple instances",
				Optional:    true,
				ForceNew:    true,
			},
			"flavor": {
				Type:        schema.TypeSet,
				Required:    true,
				MaxItems:    1,
				Description: "Flavor information",
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flavor_id": {
							Type:        schema.TypeString,
							Description: "Flavor id",
							Required:    true,
						},
					},
				},
			},
			"group": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Description: "Start instance in group",
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Description: "Group id",
							Optional:    true,
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Instance name",
				Required:    true,
				ForceNew:    true,
			},
			"ssh_key": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Existing SSH Keypair",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "SSH Keypair name",
							Required:    true,
						},
					},
				},
			},
			"ssh_key_create": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Creatting SSH Keypair",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "SSH Keypair name",
							Required:    true,
						},
						"public_key": {
							Type:        schema.TypeString,
							Description: "Group id",
							Required:    true,
						},
					},
				},
			},
			"user_data": {
				Type:        schema.TypeString,
				Description: "Configuration information or scripts to use upon launch",
				Optional:    true,
				ForceNew:    true,
			},
			"network": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Create network interfaces",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public": {
							Type:        schema.TypeBool,
							Description: "Set the new instance as public",
							Optional:    true,
						},
					},
				},
			},
			// computed
			"addresses": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Instance IP addresses",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Description: "IP address",
							Computed:    true,
						},
						"version": {
							Type:        schema.TypeInt,
							Description: "IP version",
							Computed:    true,
						},
					},
				},
			},
			"attached_volumes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: " Volumes attached to the instance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Volume Id",
							Computed:    true,
						},
					},
				},
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Description: "Flavor id",
				Computed:    true,
			},
			"flavor_name": {
				Type:        schema.TypeString,
				Description: "Flavor name",
				Computed:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "Instance id",
				Computed:    true,
			},
			"image_id": {
				Type:        schema.TypeString,
				Description: "Image id",
				Computed:    true,
			},
			"task_state": {
				Type:        schema.TypeString,
				Description: "Instance task state",
				Computed:    true,
			},
		},
	}
}

func resourceCloudProjectInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	region := d.Get("region").(string)
	params := new(CloudProjectInstanceCreateOpts)
	params.FromResource(d)

	r := &CloudProjectOperation{}

	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/instance",
		url.PathEscape(serviceName),
		url.PathEscape(region),
	)

	if err := config.OVHClient.Post(endpoint, params, r); err != nil {
		return diag.Errorf("calling %s with params %v:\n\t %q", endpoint, params, err)
	}

	err := waitForInstanceCreation(ctx, config.OVHClient, serviceName, r.Id)
	if err != nil {
		return diag.Errorf("timeout instance creation: %s", err)
	}

	endpointInstance := fmt.Sprintf("/cloud/project/%s/operation/%s",
		url.PathEscape(serviceName),
		url.PathEscape(r.Id),
	)

	err = config.OVHClient.GetWithContext(ctx, endpointInstance, r)
	if err != nil {
		return diag.Errorf("failed to get instance id: %s", err)
	}

	d.SetId(r.SubOperations[0].ResourceId)
	d.Set("region", region)

	return resourceCloudProjectInstanceRead(ctx, d, meta)
}

func waitForInstanceCreation(ctx context.Context, client *ovh.Client, serviceName, operationId string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"null", "in-progress", "created", ""},
		Target:  []string{"completed"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectOperation{}
			endpoint := fmt.Sprintf("/cloud/project/%s/operation/%s",
				url.PathEscape(serviceName),
				url.PathEscape(operationId),
			)
			err := client.GetWithContext(ctx, endpoint, res)
			if err != nil {
				return res, "", err
			}
			return res, res.Status, nil
		},
		Timeout:    360 * time.Second,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceCloudProjectInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	id := d.Id()
	region := d.Get("region").(string)
	serviceName := d.Get("service_name").(string)

	r := &CloudProjectInstanceResponse{}

	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/instance/%s",
		url.PathEscape(serviceName),
		url.PathEscape(region),
		url.PathEscape(id),
	)

	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return diag.Errorf("Error calling post %s:\n\t %q", endpoint, err)
	}

	d.Set("flavor_id", r.FlavorId)
	d.Set("flavor_name", r.FlavorName)
	d.Set("image_id", r.ImageId)
	d.Set("region", r.Region)
	d.Set("task_state", r.TaskState)
	d.Set("name", d.Get("name").(string))
	d.Set("id", r.Id)

	addresses := make([]map[string]interface{}, 0)
	if r.Addresses != nil {
		for _, add := range r.Addresses {
			address := make(map[string]interface{})
			address["ip"] = add.Ip
			address["version"] = add.Version
			addresses = append(addresses, address)
		}
	}
	d.Set("addresses", addresses)

	attachedVolumes := make([]map[string]interface{}, 0)
	if r.AttachedVolumes != nil {
		for _, att := range r.AttachedVolumes {
			attachedVolume := make(map[string]interface{})
			attachedVolume["id"] = att.Id
			attachedVolumes = append(attachedVolumes, attachedVolume)
		}
	}
	d.Set("attached_volumes", attachedVolumes)

	return nil
}

func resourceCloudProjectInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	region := d.Get("region").(string)

	id := d.Id()

	log.Printf("[DEBUG] Will delete public cloud instance for project: %s, region: %s, id: %s", serviceName, region, id)

	endpoint := fmt.Sprintf("/cloud/project/%s/instance/%s",
		url.PathEscape(serviceName),
		url.PathEscape(id),
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return diag.Errorf("Error calling post %s:\n\t %q", endpoint, err)
	}

	d.SetId("")

	log.Printf("[DEBUG] Deleted Public Cloud %s Gateway %s", serviceName, id)
	return nil
}

package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
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
				Description: "Service name of the resource representing the id of the cloud project",
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
			"availability_zone": {
				Type:        schema.TypeString,
				Description: "The availability zone where the instance will be created",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"ssh_key": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Existing SSH Key pair",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "SSH Key pair name",
							Required:    true,
						},
					},
				},
				ExactlyOneOf: []string{"ssh_key", "ssh_key_create"},
			},
			"ssh_key_create": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Add existing SSH Key pair into your Public Cloud project and link it to the instance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "SSH Key pair name",
							Required:    true,
						},
						"public_key": {
							Type:        schema.TypeString,
							Description: "SSH Public Key",
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
				Description: "Volumes attached to the instance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Volume id",
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

	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/instance",
		url.PathEscape(serviceName),
		url.PathEscape(region),
	)

	r := &CloudProjectOperationResponse{}
	if err := config.OVHClient.Post(endpoint, params, r); err != nil {
		return diag.Errorf("calling %s with params %v:\n\t %q", endpoint, params, err)
	}

	instanceID, err := waitForCloudProjectOperation(ctx, config.OVHClient, serviceName, r.Id, "instance#create")
	if err != nil {
		return diag.Errorf("timeout instance creation: %s", err)
	}

	d.SetId(instanceID)

	return resourceCloudProjectInstanceRead(ctx, d, meta)
}

func resourceCloudProjectInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	id := d.Id()
	region := d.Get("region").(string)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/instance/%s",
		url.PathEscape(serviceName),
		url.PathEscape(region),
		url.PathEscape(id),
	)

	r := &CloudProjectInstanceResponse{}
	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return diag.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
	}

	d.Set("flavor_id", r.FlavorId)
	d.Set("flavor_name", r.FlavorName)
	d.Set("image_id", r.ImageId)
	d.Set("region", r.Region)
	d.Set("task_state", r.TaskState)
	d.Set("id", r.Id)

	addresses := make([]map[string]interface{}, 0, len(r.Addresses))
	for _, add := range r.Addresses {
		address := make(map[string]interface{})
		address["ip"] = add.Ip
		address["version"] = add.Version
		addresses = append(addresses, address)
	}
	d.Set("addresses", addresses)

	attachedVolumes := make([]map[string]interface{}, 0, len(r.AttachedVolumes))
	for _, att := range r.AttachedVolumes {
		attachedVolume := make(map[string]interface{})
		attachedVolume["id"] = att.Id
		attachedVolumes = append(attachedVolumes, attachedVolume)
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
		return diag.Errorf("Error calling Delete %s:\n\t %q", endpoint, err)
	}

	d.SetId("")

	log.Printf("[DEBUG] Deleted Public Cloud %s Instance %s", serviceName, id)
	return nil
}

package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func dataSourceCloudProjectInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectInstanceRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"region": {
				Type:        schema.TypeString,
				Description: "Instance region",
				Required:    true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Description: "Instance id",
				Required:    true,
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
			"name": {
				Type:        schema.TypeString,
				Description: "Instance name",
				Computed:    true,
			},
			"image_id": {
				Type:        schema.TypeString,
				Description: "Image id",
				Computed:    true,
			},
			"ssh_key": {
				Type:        schema.TypeString,
				Description: "SSH Key pair name",
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

func dataSourceCloudProjectInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	region := d.Get("region").(string)
	instanceId := d.Get("instance_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/instance/%s",
		url.PathEscape(serviceName),
		url.PathEscape(region),
		url.PathEscape(instanceId),
	)
	var res CloudProjectInstanceResponse

	log.Printf("[DEBUG] Will read instance %s from project %s in region %s", instanceId, serviceName, region)
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	addresses := make([]map[string]interface{}, 0, len(res.Addresses))
	for _, addr := range res.Addresses {
		address := make(map[string]interface{})
		address["ip"] = addr.Ip
		address["version"] = addr.Version
		addresses = append(addresses, address)
	}

	attachedVolumes := make([]map[string]interface{}, 0, len(res.AttachedVolumes))
	for _, volume := range res.AttachedVolumes {
		attachedVolume := make(map[string]interface{})
		attachedVolume["id"] = volume.Id
		attachedVolumes = append(attachedVolumes, attachedVolume)
	}

	d.Set("addresses", addresses)
	d.Set("flavor_id", res.FlavorId)
	d.Set("flavor_name", res.FlavorName)
	d.SetId(res.Id)
	d.Set("image_id", res.ImageId)
	d.Set("instance_id", res.Id)
	d.Set("name", res.Name)
	d.Set("ssh_key", res.SshKey)
	d.Set("task_state", res.TaskState)
	d.Set("attached_volumes", attachedVolumes)

	return nil
}

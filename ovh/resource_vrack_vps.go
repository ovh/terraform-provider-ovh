package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceVrackVps() *schema.Resource {
	return &schema.Resource{
		Create: resourceVrackVpsCreate,
		Read:   resourceVrackVpsRead,
		Delete: resourceVrackVpsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVrackVpsImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_VRACK_SERVICE", nil),
				Description: "Service name of the vrack resource.",
			},
			"vps_service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Service name of the VPS to attach to the vRack.",
			},
		},
	}
}

func resourceVrackVpsImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not SERVICE_NAME/VPS_SERVICE_NAME formatted")
	}
	serviceName := splitId[0]
	vpsServiceName := splitId[1]
	d.SetId(fmt.Sprintf("vrack_%s-vps_%s", serviceName, vpsServiceName))
	d.Set("service_name", serviceName)
	d.Set("vps_service_name", vpsServiceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceVrackVpsCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)

	opts := (&VrackVpsCreateOpts{}).FromResource(d)
	task := &VrackTask{}

	endpoint := fmt.Sprintf("/vrack/%s/vps", url.PathEscape(serviceName))

	if err := config.OVHClient.Post(endpoint, opts, task); err != nil {
		return fmt.Errorf("Error calling POST %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	// Only wait for the task if the API actually returned one with a populated id/serviceName.
	// Some vrack endpoints return a void body for this call.
	if task != nil && task.Id != 0 && task.ServiceName != "" {
		if err := waitForVrackTask(task, config.OVHClient); err != nil {
			return fmt.Errorf("Error waiting for vrack (%s) to attach VPS %s: %s", serviceName, opts.Vps, err)
		}
	}

	d.SetId(fmt.Sprintf("vrack_%s-vps_%s", serviceName, opts.Vps))

	return resourceVrackVpsRead(d, meta)
}

func resourceVrackVpsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	vv := &VrackVps{}
	serviceName := d.Get("service_name").(string)
	vpsServiceName := d.Get("vps_service_name").(string)

	endpoint := fmt.Sprintf("/vrack/%s/vps/%s",
		url.PathEscape(serviceName),
		url.PathEscape(vpsServiceName),
	)

	if err := config.OVHClient.Get(endpoint, vv); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("service_name", vv.Vrack)
	d.Set("vps_service_name", vv.Vps)

	return nil
}

func resourceVrackVpsDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	vpsServiceName := d.Get("vps_service_name").(string)

	task := &VrackTask{}
	endpoint := fmt.Sprintf("/vrack/%s/vps/%s",
		url.PathEscape(serviceName),
		url.PathEscape(vpsServiceName),
	)

	if err := config.OVHClient.Delete(endpoint, task); err != nil {
		return fmt.Errorf("Error calling DELETE %s with %s/%s:\n\t %q", endpoint, serviceName, vpsServiceName, err)
	}

	if task != nil && task.Id != 0 && task.ServiceName != "" {
		if err := waitForVrackTask(task, config.OVHClient); err != nil {
			return fmt.Errorf("Error waiting for vrack (%s) to detach VPS (%s): %s", serviceName, vpsServiceName, err)
		}
	}

	d.SetId("")
	return nil
}

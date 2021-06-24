package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceVrackIp() *schema.Resource {
	return &schema.Resource{
		Create: resourceVrackIpCreate,
		Read:   resourceVrackIpRead,
		Delete: resourceVrackIpDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVrackIpImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of your vrack",
			},
			"block": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Your IP block.",
			},

			// computed
			"gateway": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Your gateway",
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Your IP block",
			},
			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Where you want your block announced on the network",
			},
		},
	}
}

func resourceVrackIpImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, ",", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not SERVICE_NAME,IPBlock formatted")
	}
	serviceName := splitId[0]
	block := splitId[1]
	d.SetId(fmt.Sprintf("vrack_%s-block_%s", serviceName, block))
	d.Set("service_name", serviceName)
	d.Set("block", block)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceVrackIpCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)

	opts := (&VrackIpCreateOpts{}).FromResource(d)
	task := &VrackTask{}

	endpoint := fmt.Sprintf("/vrack/%s/ip", serviceName)

	if err := config.OVHClient.Post(endpoint, opts, task); err != nil {
		return fmt.Errorf("Error calling POST %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	if err := waitForVrackTask(task, config.OVHClient); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to attach ip %v: %s", serviceName, opts, err)
	}

	//set id
	d.SetId(fmt.Sprintf("vrack_%s-dedicatedserver_%s", serviceName, opts.Block))

	return resourceVrackIpRead(d, meta)
}

func resourceVrackIpRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	block := d.Get("block").(string)

	ip := &VrackIp{}
	endpoint := fmt.Sprintf("/vrack/%s/ip/%s",
		url.PathEscape(serviceName),
		url.PathEscape(block),
	)

	if err := config.OVHClient.Get(endpoint, ip); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("block", ip.Ip)

	// set resource attributes
	for k, v := range ip.ToMap() {
		d.Set(k, v)
	}

	return nil
}

func resourceVrackIpDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	block := d.Get("block").(string)
	task := &VrackTask{}

	endpoint := fmt.Sprintf("/vrack/%s/ip/%s",
		url.PathEscape(serviceName),
		url.PathEscape(block),
	)

	if err := config.OVHClient.Delete(endpoint, task); err != nil {
		return fmt.Errorf("Error calling DELETE %s with %s/%s:\n\t %q", endpoint, serviceName, block, err)
	}

	if err := waitForVrackTask(task, config.OVHClient); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to detach ip (%s): %s", serviceName, block, err)
	}

	d.SetId("")
	return nil
}

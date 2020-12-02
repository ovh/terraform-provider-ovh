package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceVrackIpLoadbalancing() *schema.Resource {
	return &schema.Resource{
		Create: resourceVrackIpLoadbalancingCreate,
		Read:   resourceVrackIpLoadbalancingRead,
		Delete: resourceVrackIpLoadbalancingDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVrackIpLoadbalancingImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The internal name of your vrack",
				Required:    true,
				ForceNew:    true,
			},
			"ip_loadbalancing": {
				Type:        schema.TypeString,
				Description: "Your ipLoadbalancing",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceVrackIpLoadbalancingImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not service_name/iploadbalancing formatted")
	}
	serviceName := splitId[0]
	ipLoadbalancing := splitId[1]
	d.SetId(fmt.Sprintf("%s-%s", serviceName, ipLoadbalancing))
	d.Set("service_name", serviceName)
	d.Set("ip_loadbalancing", ipLoadbalancing)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceVrackIpLoadbalancingCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	opts := (&VrackIpLoadbalancingCreateOpts{}).FromResource(d)
	task := &VrackTask{}

	endpoint := fmt.Sprintf("/vrack/%s/ipLoadbalancing", serviceName)

	if err := config.OVHClient.Post(endpoint, opts, task); err != nil {
		return fmt.Errorf("Error calling POST %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	if err := waitForVrackTask(task, config.OVHClient); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to attach dedicated server %v: %s", serviceName, opts, err)
	}

	//set id
	d.SetId(fmt.Sprintf("%s-%s", serviceName, opts.IpLoadbalancing))

	return resourceVrackIpLoadbalancingRead(d, meta)
}

func resourceVrackIpLoadbalancingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	vds := &VrackIpLoadbalancing{}

	serviceName := d.Get("service_name").(string)
	ipLoadbalancing := d.Get("ip_loadbalancing").(string)

	endpoint := fmt.Sprintf("/vrack/%s/ipLoadbalancing/%s",
		url.PathEscape(serviceName),
		url.PathEscape(ipLoadbalancing),
	)

	if err := config.OVHClient.Get(endpoint, vds); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("service_name", vds.Vrack)
	d.Set("ip_loadbalancing", vds.IpLoadbalancing)

	return nil
}

func resourceVrackIpLoadbalancingDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	ipLoadbalancing := d.Get("ip_loadbalancing").(string)

	task := &VrackTask{}
	endpoint := fmt.Sprintf("/vrack/%s/ipLoadbalancing/%s",
		url.PathEscape(serviceName),
		url.PathEscape(ipLoadbalancing),
	)

	if err := config.OVHClient.Delete(endpoint, task); err != nil {
		return fmt.Errorf("Error calling DELETE %s with %s/%s:\n\t %q", endpoint, serviceName, ipLoadbalancing, err)
	}

	if err := waitForVrackTask(task, config.OVHClient); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to detach dedicated server (%s): %s", serviceName, ipLoadbalancing, err)
	}

	d.SetId("")
	return nil
}

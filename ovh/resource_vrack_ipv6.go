package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceVrackIpV6() *schema.Resource {
	return &schema.Resource{
		Create: resourceVrackIpv6Create,
		Read:   resourceVrackIpv6Read,
		Delete: resourceVrackIpv6Delete,
		Importer: &schema.ResourceImporter{
			State: resourceVrackIpv6ImportState,
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
				Description: "IPv6 CIDR notation (e.g., 2001:41d0::/128)",
				ValidateFunc: func(i interface{}, _ string) ([]string, []error) {
					if err := helpers.ValidateIpBlock(i.(string)); err != nil {
						return nil, []error{err}
					}
					return nil, nil
				},
			},
		},
	}
}

func resourceVrackIpv6ImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.Split(givenId, ",")
	if len(splitId) != 2 {
		return nil, fmt.Errorf("import ID is not SERVICE_NAME,IPv6-block formatted")
	}
	serviceName := splitId[0]
	block := splitId[1]

	d.SetId(fmt.Sprintf("vrack_%s-block_%s", serviceName, block))
	d.Set("service_name", serviceName)
	d.Set("block", block)

	return []*schema.ResourceData{d}, nil
}

func resourceVrackIpv6Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)

	opts := (&VrackIpCreateOpts{}).FromResource(d)
	task := VrackTask{}

	endpoint := fmt.Sprintf("/vrack/%s/ipv6", url.PathEscape(serviceName))
	if err := config.OVHClient.Post(endpoint, opts, &task); err != nil {
		return fmt.Errorf("error calling POST %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	if err := waitForVrackTask(&task, config.OVHClient); err != nil {
		return fmt.Errorf("error waiting for vrack (%s) to attach ipv6 %v: %s", serviceName, opts, err)
	}

	d.SetId(fmt.Sprintf("vrack_%s-block_%s", serviceName, opts.Block))

	return resourceVrackIpv6Read(d, meta)
}

func resourceVrackIpv6Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	block := d.Get("block").(string)

	endpoint := fmt.Sprintf("/vrack/%s/ipv6/%s",
		url.PathEscape(serviceName),
		url.PathEscape(block),
	)

	if err := config.OVHClient.Get(endpoint, nil); err != nil {
		return fmt.Errorf("failed to get vrack-ipv6 link: %w", err)
	}

	return nil
}

func resourceVrackIpv6Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	block := d.Get("block").(string)
	task := VrackTask{}

	endpoint := fmt.Sprintf("/vrack/%s/ipv6/%s",
		url.PathEscape(serviceName),
		url.PathEscape(block),
	)

	if err := config.OVHClient.Delete(endpoint, &task); err != nil {
		return fmt.Errorf("error calling DELETE %s with %s/%s:\n\t %q", endpoint, serviceName, block, err)
	}

	if err := waitForVrackTask(&task, config.OVHClient); err != nil {
		return fmt.Errorf("error waiting for vrack (%s) to detach ip (%s): %s", serviceName, block, err)
	}

	d.SetId("")

	return nil
}

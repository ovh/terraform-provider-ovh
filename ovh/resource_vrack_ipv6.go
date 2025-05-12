package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceVrackIpV6() *schema.Resource {
	return &schema.Resource{
		Create: resourceVrackIpv6Create,
		Update: resourceVrackIpv6Update,
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
			"bridged_subrange": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Computed:    true,
				Optional:    true,
				Description: "Subrange bridged into your vrack",
				ForceNew:    false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subrange": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv6 CIDR notation (e.g., 2001:41d0::/128)",
						},
						"gateway": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Your gateway",
						},
						"slaac": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Slaac status",
							ValidateFunc: helpers.ValidateEnum([]string{"disabled", "enabled"}),
						},
					},
				},
			},
			// computed
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Where your block announced on the network",
			},
			"ipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPv6 block announced on the network",
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

	if err := setBridgedSubrangeState(d, meta, serviceName, block); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func resourceVrackIpv6Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)

	opts := (&VrackIpV6CreateOpts{}).FromResource(d)
	task := VrackTask{}

	endpoint := fmt.Sprintf("/vrack/%s/ipv6", url.PathEscape(serviceName))
	if err := config.OVHClient.Post(endpoint, opts, &task); err != nil {
		return fmt.Errorf("error calling POST %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	if err := waitForVrackTask(&task, config.OVHClient); err != nil {
		return fmt.Errorf("error waiting for vrack (%s) to attach ipv6 %v: %s", serviceName, opts, err)
	}

	optSlaac := (&VrackIPv6BridgedSubrangeSlaacUpdateOpts{}).FromResource(d)
	if optSlaac.Slaac == "enabled" {
		var task VrackTask
		endpoint := fmt.Sprintf("/vrack/%s/ipv6/%s/bridgedSubrange",
			url.PathEscape(serviceName),
			url.PathEscape(opts.Block),
		)

		log.Printf("[DEBUG] Get the subrange bridged into your vrack")
		var bridgedSubranges []string
		if err := config.OVHClient.Get(endpoint, &bridgedSubranges); err != nil {
			return fmt.Errorf("error calling Get %s: %w", endpoint, err)
		}
		if len(bridgedSubranges) != 1 {
			return fmt.Errorf("error getting bridgeSubrange: exactly one should be found")
		}
		bridgedSubrange := bridgedSubranges[0]

		log.Printf("[DEBUG] Will set bridged subrange %s Slaac to %s", bridgedSubrange, optSlaac.Slaac)
		endpoint = fmt.Sprintf("/vrack/%s/ipv6/%s/bridgedSubrange/%s",
			url.PathEscape(serviceName),
			url.PathEscape(opts.Block),
			url.PathEscape(bridgedSubrange),
		)
		if err := config.OVHClient.Put(endpoint, optSlaac, &task); err != nil {
			return fmt.Errorf("error calling Put %s: %q", endpoint, err)
		}

		if err := waitForVrackTask(&task, config.OVHClient); err != nil {
			return fmt.Errorf("error waiting for vrack (%s): %s", serviceName, err)
		}
	}
	d.SetId(fmt.Sprintf("vrack_%s-block_%s", serviceName, opts.Block))

	return resourceVrackIpv6Read(d, meta)
}

func resourceVrackIpv6Update(d *schema.ResourceData, meta interface{}) error {
	var task VrackTask
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	block := d.Get("block").(string)

	log.Printf("[DEBUG] Get the subrange bridged into your vrack")
	endpoint := fmt.Sprintf("/vrack/%s/ipv6/%s/bridgedSubrange",
		url.PathEscape(serviceName),
		url.PathEscape(block),
	)

	var bridgedSubranges []string
	if err := config.OVHClient.Get(endpoint, &bridgedSubranges); err != nil {
		return fmt.Errorf("error calling Get %s: %w", endpoint, err)
	}
	if len(bridgedSubranges) != 1 {
		return fmt.Errorf("error getting bridgeSubrange: exactly one should be found")
	}
	bridgedSubrange := bridgedSubranges[0]

	opts := (&VrackIPv6BridgedSubrangeSlaacUpdateOpts{}).FromResource(d)
	log.Printf("[DEBUG] Will update bridge subrange %s slaac to: %s", bridgedSubrange, opts.Slaac)
	endpoint = fmt.Sprintf("/vrack/%s/ipv6/%s/bridgedSubrange/%s",
		url.PathEscape(serviceName),
		url.PathEscape(block),
		url.PathEscape(bridgedSubrange),
	)

	if err := config.OVHClient.Put(endpoint, opts, &task); err != nil {
		return fmt.Errorf("error calling Put %s: %q", endpoint, err)
	}

	if err := waitForVrackTask(&task, config.OVHClient); err != nil {
		return fmt.Errorf("error waiting for vrack (%s): %s", serviceName, err)
	}

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

	ipv6 := &VrackIpV6{}
	if err := config.OVHClient.Get(endpoint, ipv6); err != nil {
		return fmt.Errorf("failed to get vrack-ipv6 link: %w", err)
	}

	// set resource attributes
	for k, v := range ipv6.ToMap() {
		d.Set(k, v)
	}

	return setBridgedSubrangeState(d, meta, serviceName, block)
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

// setBridgedSubrangeState set the values of the BridgedSubrange nested resource in the state.
func setBridgedSubrangeState(d *schema.ResourceData, meta interface{}, serviceName, block string) error {
	config := meta.(*Config)

	endpoint := fmt.Sprintf("/vrack/%s/ipv6/%s/bridgedSubrange",
		url.PathEscape(serviceName),
		url.PathEscape(block),
	)

	log.Printf("[DEBUG] Get the subrange bridged into your vrack")
	var bridgedSubranges []string
	if err := config.OVHClient.Get(endpoint, &bridgedSubranges); err != nil {
		return fmt.Errorf("error calling Get %s: %w", endpoint, err)
	}

	if len(bridgedSubranges) != 1 {
		return fmt.Errorf("error getting bridgeSubrange: exactly one should be found")
	}
	endpoint = fmt.Sprintf("/vrack/%s/ipv6/%s/bridgedSubrange/%s",
		url.PathEscape(serviceName),
		url.PathEscape(block),
		url.PathEscape(bridgedSubranges[0]),
	)

	var bridgedSubrange VrackIPv6BridgedSubrange
	if err := config.OVHClient.Get(endpoint, &bridgedSubrange); err != nil {
		return fmt.Errorf("error calling Get %s: %w", endpoint, err)
	}

	d.Set("bridged_subrange", bridgedSubrange.ToMap())

	return nil
}

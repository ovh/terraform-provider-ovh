package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceVPSIpReverse() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPSIpReverseCreate,
		Read:   resourceVPSIpReverseRead,
		Update: resourceVPSIpReverseUpdate,
		Delete: resourceVPSIpReverseDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVPSIpReverseImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"reverse": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Computed - geolocation is read-only (set at order time, not via PUT)
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"geolocation": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVPSIpReverseImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	parts := strings.SplitN(givenId, "|", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("Import Id must be service_name|ip_address (got %q)", givenId)
	}
	d.Set("service_name", parts[0])
	d.Set("ip_address", parts[1])
	d.SetId(parts[0] + "|" + parts[1])
	return []*schema.ResourceData{d}, nil
}

func vpsIpReverseEndpoint(serviceName, ipAddress string) string {
	return fmt.Sprintf(
		"/vps/%s/ips/%s",
		url.PathEscape(serviceName),
		url.PathEscape(ipAddress),
	)
}

func resourceVPSIpReversePut(d *schema.ResourceData, meta interface{}, reverse string) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	ipAddress := d.Get("ip_address").(string)

	opts := VPSIpReverseUpdateOpts{Reverse: reverse}
	endpoint := vpsIpReverseEndpoint(serviceName, ipAddress)
	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return fmt.Errorf("Error calling PUT %s: %q", endpoint, err)
	}
	return nil
}

func resourceVPSIpReverseCreate(d *schema.ResourceData, meta interface{}) error {
	serviceName := d.Get("service_name").(string)
	ipAddress := d.Get("ip_address").(string)
	reverse := d.Get("reverse").(string)

	if err := resourceVPSIpReversePut(d, meta, reverse); err != nil {
		return err
	}
	d.SetId(serviceName + "|" + ipAddress)
	return resourceVPSIpReverseRead(d, meta)
}

func resourceVPSIpReverseUpdate(d *schema.ResourceData, meta interface{}) error {
	reverse := d.Get("reverse").(string)
	if err := resourceVPSIpReversePut(d, meta, reverse); err != nil {
		return err
	}
	return resourceVPSIpReverseRead(d, meta)
}

func resourceVPSIpReverseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	ipAddress := d.Get("ip_address").(string)
	endpoint := vpsIpReverseEndpoint(serviceName, ipAddress)

	ip := &VPSIp{}
	if err := config.OVHClient.Get(endpoint, ip); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("version", ip.Version)
	d.Set("type", ip.Type)
	if ip.Gateway != nil {
		d.Set("gateway", *ip.Gateway)
	}
	d.Set("geolocation", ip.Geolocation)
	if ip.Reverse != nil {
		d.Set("reverse", *ip.Reverse)
	} else {
		d.Set("reverse", "")
	}
	return nil
}

// resourceVPSIpReverseDelete clears the reverse DNS record by PUTting an empty
// `reverse` value. It deliberately does NOT issue DELETE /vps/{sn}/ips/{ip} —
// that endpoint *releases the additional IP*, which is a destructive billing
// operation that must be handled via the cart termination flow.
func resourceVPSIpReverseDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Clearing reverse DNS for vps=%s ip=%s",
		d.Get("service_name").(string), d.Get("ip_address").(string))
	if err := resourceVPSIpReversePut(d, meta, ""); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

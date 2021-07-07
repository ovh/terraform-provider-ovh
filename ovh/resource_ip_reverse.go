package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceIpReverse() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpReverseCreate,
		Read:   resourceIpReverseRead,
		Delete: resourceIpReverseDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIpReverseImportState,
		},

		Schema: map[string]*schema.Schema{
			"ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateIpBlock(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},

			"ip_reverse": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateIp(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},

			"reverse": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceIpReverseImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, ":", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not ip:ip_reverse formatted")
	}
	ip := splitId[0]
	ipReverse := splitId[1]
	d.SetId(ipReverse)
	d.Set("ip", ip)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceIpReverseCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Create the new reverse
	ip := d.Get("ip").(string)
	opts := (&IpReverseCreateOpts{}).FromResource(d)
	res := &IpReverse{}

	err := config.OVHClient.Post(
		fmt.Sprintf("/ip/%s/reverse", url.PathEscape(ip)),
		opts,
		&res,
	)
	if err != nil {
		return fmt.Errorf("Failed to create OVH IP Reverse: %s", err)
	}

	d.SetId(res.IpReverse)

	return resourceIpReverseRead(d, meta)
}

func resourceIpReverseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ip := d.Get("ip").(string)

	res := &IpReverse{}
	endpoint := fmt.Sprintf(
		"/ip/%s/reverse/%s",
		url.PathEscape(ip),
		url.PathEscape(d.Id()),
	)

	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	for k, v := range res.ToMap() {
		d.Set(k, v)
	}

	return nil
}

func resourceIpReverseDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	log.Printf("[INFO] Deleting OVH IP Reverse: %s->%s", d.Get("reverse").(string), d.Get("ip_reverse").(string))
	ip := d.Get("ip").(string)
	endpoint := fmt.Sprintf(
		"/ip/%s/reverse/%s",
		url.PathEscape(ip),
		url.PathEscape(d.Id()),
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId("")
	return nil
}

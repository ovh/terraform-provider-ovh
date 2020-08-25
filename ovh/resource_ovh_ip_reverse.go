package ovh

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-ovh/ovh/helpers"

	"github.com/ovh/go-ovh/ovh"
)

type OvhIpReverse struct {
	IpReverse string `json:"ipReverse"`
	Reverse   string `json:"reverse"`
}

func resourceOvhIpReverse() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvhIpReverseCreate,
		Read:   resourceOvhIpReverseRead,
		Update: resourceOvhIpReverseUpdate,
		Delete: resourceOvhIpReverseDelete,

		Schema: map[string]*schema.Schema{
			"ip": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateIpBlock(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"ipreverse": {
				Type:     schema.TypeString,
				Optional: true,
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
				Required: true,
			},
		},
	}
}

func resourceOvhIpReverseCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Config)

	// Create the new reverse
	newIp := d.Get("ip").(string)
	newReverse := &OvhIpReverse{
		Reverse: d.Get("reverse").(string),
	}

	newIpReverse, ok := d.GetOk("ipreverse")
	if !ok || newIpReverse == "" {
		ipAddr, ipNet, _ := net.ParseCIDR(newIp)
		prefixSize, _ := ipNet.Mask.Size()

		if ipAddr.To4() != nil && prefixSize != 32 {
			return fmt.Errorf("ipreverse must be set if ip (%s) is not a /32", newIp)
		} else if ipAddr.To4() == nil && prefixSize != 128 {
			return fmt.Errorf("ipreverse must be set if ip (%s) is not a /128", newIp)
		}

		newIpReverse = ipAddr.String()
		d.Set("ipreverse", newIpReverse)
	}

	newReverse.IpReverse = newIpReverse.(string)

	log.Printf("[DEBUG] OVH IP Reverse create configuration: %#v", newReverse)

	resultReverse := OvhIpReverse{}

	err := provider.OVHClient.Post(
		fmt.Sprintf("/ip/%s/reverse", strings.Replace(newIp, "/", "%2F", 1)),
		newReverse,
		&resultReverse,
	)
	if err != nil {
		return fmt.Errorf("Failed to create OVH IP Reverse: %s", err)
	}

	d.SetId(fmt.Sprintf("%s_%s", newIp, resultReverse.IpReverse))

	return resourceOvhIpReverseRead(d, meta)
}

func resourceOvhIpReverseRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Config)

	reverse := OvhIpReverse{}
	err := provider.OVHClient.Get(
		fmt.Sprintf("/ip/%s/reverse/%s", strings.Replace(d.Get("ip").(string), "/", "%2F", 1), d.Get("ipreverse").(string)),
		&reverse,
	)

	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("ipreverse", reverse.IpReverse)
	d.Set("reverse", reverse.Reverse)

	return nil
}

func resourceOvhIpReverseUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Config)

	reverse := OvhIpReverse{}

	if attr, ok := d.GetOk("ipreverse"); ok {
		reverse.IpReverse = attr.(string)
	}
	if attr, ok := d.GetOk("reverse"); ok {
		reverse.Reverse = attr.(string)
	}

	log.Printf("[DEBUG] OVH IP Reverse update configuration: %#v", reverse)

	err := provider.OVHClient.Post(
		fmt.Sprintf("/ip/%s/reverse", strings.Replace(d.Get("ip").(string), "/", "%2F", 1)),
		reverse,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Failed to update OVH IP Reverse: %s", err)
	}

	return resourceOvhIpReverseRead(d, meta)
}

func resourceOvhIpReverseDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Config)

	log.Printf("[INFO] Deleting OVH IP Reverse: %s->%s", d.Get("reverse").(string), d.Get("ipreverse").(string))

	err := provider.OVHClient.Delete(
		fmt.Sprintf("/ip/%s/reverse/%s", strings.Replace(d.Get("ip").(string), "/", "%2F", 1), d.Get("ipreverse").(string)),
		nil,
	)

	if err != nil {
		return fmt.Errorf("Error deleting OVH IP Reverse: %s", err)
	}

	return nil
}

func resourceOvhIpReverseExists(ip, ipreverse string, c *ovh.Client) error {
	reverse := OvhIpReverse{}
	endpoint := fmt.Sprintf("/ip/%s/reverse/%s", strings.Replace(ip, "/", "%2F", 1), ipreverse)

	err := c.Get(endpoint, &reverse)
	if err != nil {
		return fmt.Errorf("calling %s:\n\t %q", endpoint, err)
	}
	log.Printf("[DEBUG] Read IP reverse: %s", reverse)

	return nil
}

package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceIpMove() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpMoveCreate,
		Read:   resourceIpMoveRead,
		Delete: resourceIpMoveDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
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
			"nexthop": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "",
			},
			"to": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceIpMoveCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Create the new move
	ip := d.Get("ip").(string)
	opts := (&IpMoveCreateOpts{}).FromResource(d)
	res := &IpTask{}

	err := config.OVHClient.Post(
		fmt.Sprintf("/ip/%s/move", url.PathEscape(ip)), opts, &res,
	)
	if err != nil {
		return fmt.Errorf("Failed to create OVH IP Move: %s", err)
	}

	return nil
}

func resourceIpMoveRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ip := d.Get("ip").(string)

	res := &IpDestinations{}
	endpoint := fmt.Sprintf("/ip/%s/move", url.PathEscape(ip))

	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	for k, v := range res.ToMap() {
		d.Set(k, v)
	}

	return nil
}

func resourceIpMoveDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

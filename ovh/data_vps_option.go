package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

// vpsOptionEnum is the canonical list of vps.VpsOptionEnum values
// accepted by the OVH API for /vps/{serviceName}/option/{option}.
var vpsOptionEnum = []string{
	"additionalDisk",
	"automatedBackup",
	"cpanel",
	"ftpbackup",
	"plesk",
	"snapshot",
	"veeam",
	"windows",
}

func dataSourceVPSOption() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSOptionRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your VPS.",
			},
			"option": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The option name. One of: additionalDisk, automatedBackup, cpanel, ftpbackup, plesk, snapshot, veeam, windows.",
				ValidateFunc: helpers.ValidateEnum(vpsOptionEnum),
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subscription state of the option (e.g. released, subscribed).",
			},
		},
	}
}

func dataSourceVPSOptionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	option := d.Get("option").(string)

	endpoint := fmt.Sprintf("/vps/%s/option/%s",
		url.PathEscape(serviceName),
		url.PathEscape(option),
	)
	opt := &vpsOption{}
	if err := config.OVHClient.Get(endpoint, opt); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", serviceName, option))
	d.Set("state", opt.State)
	return nil
}

package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVpsReboot() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpsRebootCreate,
		Read:   resourceVpsPowerActionReadNoop,
		Delete: resourceVpsPowerActionDeleteNoop,

		Schema: vpsPowerActionSchema(nil),
	}
}

func resourceVpsRebootCreate(d *schema.ResourceData, meta interface{}) error {
	if _, err := runVpsPowerAction(d, meta, "reboot"); err != nil {
		return err
	}
	return resourceVpsPowerActionReadNoop(d, meta)
}

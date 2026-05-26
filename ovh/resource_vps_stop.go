package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVpsStop() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpsStopCreate,
		Read:   resourceVpsPowerActionReadNoop,
		Delete: resourceVpsPowerActionDeleteNoop,

		Schema: vpsPowerActionSchema(nil),
	}
}

func resourceVpsStopCreate(d *schema.ResourceData, meta interface{}) error {
	if _, err := runVpsPowerAction(d, meta, "stop"); err != nil {
		return err
	}
	return resourceVpsPowerActionReadNoop(d, meta)
}

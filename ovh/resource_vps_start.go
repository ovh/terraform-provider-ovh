package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVpsStart() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpsStartCreate,
		Read:   resourceVpsPowerActionReadNoop,
		Delete: resourceVpsPowerActionDeleteNoop,

		Schema: vpsPowerActionSchema(nil),
	}
}

func resourceVpsStartCreate(d *schema.ResourceData, meta interface{}) error {
	if _, err := runVpsPowerAction(d, meta, "start"); err != nil {
		return err
	}
	return resourceVpsPowerActionReadNoop(d, meta)
}

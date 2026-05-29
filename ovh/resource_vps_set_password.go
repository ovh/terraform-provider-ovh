package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVpsSetPassword() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpsSetPasswordCreate,
		Read:   resourceVpsPowerActionReadNoop,
		Delete: resourceVpsPowerActionDeleteNoop,

		Schema: vpsPowerActionSchema(map[string]*schema.Schema{
			// The OVH API does not expose whether the install was performed
			// with `do_not_send_password`, so this attribute defaults to true.
			// Document: set it to false in your Terraform config if your VPS
			// was installed with `do_not_send_password = true` and you handle
			// the new password out-of-band.
			"password_sent_via_email": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether OVH emailed the new root password. True unless the VPS was installed with do_not_send_password=true (the API does not report this; the value is informational).",
			},
		}),
	}
}

func resourceVpsSetPasswordCreate(d *schema.ResourceData, meta interface{}) error {
	if _, err := runVpsPowerAction(d, meta, "setPassword"); err != nil {
		return err
	}
	// Default to true unless the user explicitly opted out via config.
	if _, ok := d.GetOk("password_sent_via_email"); !ok {
		_ = d.Set("password_sent_via_email", true)
	}
	return resourceVpsPowerActionReadNoop(d, meta)
}

package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceMeIpxeScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceMeIpxeScriptCreate,
		Read:   resourceMeIpxeScriptRead,
		Delete: resourceMeIpxeScriptDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of your script",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "For documentation purpose only. This attribute is not passed to the OVH API as it cannot be retrieved back. Instead a fake description ('$name auto description') is passed at creation time.",
			},
			"script": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Content of your IPXE script",
			},
		},
	}
}

// Common function with the datasource
func resourceMeIpxeScriptRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	r := &MeIpxeScriptResponse{}

	endpoint := fmt.Sprintf("/me/ipxeScript/%s", url.PathEscape(d.Id()))

	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("name", r.Name)
	d.Set("script", r.Script)

	return nil
}

func resourceMeIpxeScriptCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Get("name").(string)
	script := d.Get("script").(string)

	params := &MeIpxeScriptCreateOpts{
		Description: fmt.Sprintf("%s auto description", name),
		Name:        name,
		Script:      script,
	}

	response := &MeIpxeScriptResponse{}

	log.Printf("[DEBUG] Will create IpxeScript: %s", params)

	err := config.OVHClient.Post("/me/ipxeScript", params, response)
	if err != nil {
		return fmt.Errorf("Error creating IpxeScript with params %s:\n\t %q", params, err)
	}

	d.SetId(response.Name)

	return resourceMeIpxeScriptRead(d, meta)
}

func resourceMeIpxeScriptDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	err := config.OVHClient.Delete(
		fmt.Sprintf("/me/ipxeScript/%s", url.PathEscape(d.Id())),
		nil,
	)
	if err != nil {
		return fmt.Errorf("Unable to delete IpxeScript named %s:\n\t %q", d.Id(), err)
	}

	log.Printf("[DEBUG] Deleted IpxeScript %s", d.Id())
	d.SetId("")
	return nil
}

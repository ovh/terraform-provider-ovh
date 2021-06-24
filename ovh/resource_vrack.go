package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceVrack() *schema.Resource {
	return &schema.Resource{
		Create: resourceVrackCreate,
		Update: resourceVrackUpdate,
		Read:   resourceVrackRead,
		Delete: resourceVrackDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: resourceVrackSchema(),
	}
}

func resourceVrackSchema() map[string]*schema.Schema {
	schema := map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Description: "yourvrackdescription",
			Optional:    true,
			Computed:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "yourvrackname",
			Optional:    true,
			Computed:    true,
		},

		// computed
		"service_name": {
			Type:        schema.TypeString,
			Description: "The internal name of your vrack",
			Computed:    true,
		},
	}

	for k, v := range genericOrderSchema(false) {
		schema[k] = v
	}

	return schema
}

func resourceVrackCreate(d *schema.ResourceData, meta interface{}) error {
	if err := orderCreate(d, meta, "vrack"); err != nil {
		return fmt.Errorf("Could not order vrack: %q", err)
	}

	return resourceVrackUpdate(d, meta)
}

func resourceVrackUpdate(d *schema.ResourceData, meta interface{}) error {
	_, details, err := orderRead(d, meta)
	if err != nil {
		return fmt.Errorf("Could not read vrack order: %q", err)
	}

	config := meta.(*Config)
	serviceName := details[0].Domain

	log.Printf("[DEBUG] Will update vrack: %s", serviceName)
	opts := (&VrackUpdateOpts{}).FromResource(d)
	endpoint := fmt.Sprintf("/vrack/%s", serviceName)
	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return fmt.Errorf("calling Put %s: %q", endpoint, err)
	}

	return resourceVrackRead(d, meta)
}

func resourceVrackRead(d *schema.ResourceData, meta interface{}) error {
	_, details, err := orderRead(d, meta)
	if err != nil {
		return fmt.Errorf("Could not read vrack order: %q", err)
	}

	config := meta.(*Config)
	serviceName := details[0].Domain

	log.Printf("[DEBUG] Will read vrack: %s", serviceName)
	r := &Vrack{}
	endpoint := fmt.Sprintf("/vrack/%s", serviceName)
	if err := config.OVHClient.Get(endpoint, &r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("service_name", serviceName)

	// set resource attributes
	for k, v := range r.ToMap() {
		d.Set(k, v)
	}

	return nil
}

func resourceVrackDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	serviceName := d.Get("service_name").(string)
	log.Printf(
		`[WARN] The API doesn't provide any delete mechanism for VRACK.
The vrack %s (order %s) will be forgotten without being effectively deleted.`,
		serviceName,
		id,
	)
	d.SetId("")
	return nil
}

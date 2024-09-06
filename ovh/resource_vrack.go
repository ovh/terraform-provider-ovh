package ovh

import (
	"fmt"
	"log"
	"net/url"

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
		"urn": {
			Type:     schema.TypeString,
			Computed: true,
		},
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
	// Order vRack and wait for it to be delivered
	if err := orderCreateFromResource(d, meta, "vrack", true); err != nil {
		return fmt.Errorf("could not order vrack: %q", err)
	}

	// Retrieve serviceName from order ID
	_, details, err := orderReadInResource(d, meta)
	if err != nil {
		return fmt.Errorf("could not read vrack order: %q", err)
	}
	serviceName := details[0].Domain

	d.SetId(serviceName)
	d.Set("service_name", serviceName)

	return resourceVrackUpdate(d, meta)
}

func resourceVrackUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Id()

	log.Printf("[DEBUG] Will update vrack: %s", serviceName)
	opts := (&VrackUpdateOpts{}).FromResource(d)
	endpoint := fmt.Sprintf("/vrack/%s", url.PathEscape(serviceName))
	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return fmt.Errorf("calling Put %s: %q", endpoint, err)
	}

	return resourceVrackRead(d, meta)
}

func resourceVrackRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Id()

	d.Set("service_name", serviceName)

	log.Printf("[DEBUG] Will read vrack: %s", serviceName)
	r := &Vrack{}
	endpoint := fmt.Sprintf("/vrack/%s", url.PathEscape(serviceName))
	if err := config.OVHClient.Get(endpoint, &r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	// Set resource attributes
	for k, v := range r.ToMap() {
		d.Set(k, v)
	}

	// Retrieve order information
	serviceObj, err := serviceFromServiceName(config.OVHClient, "vrack", serviceName)
	if err != nil {
		return fmt.Errorf("failed to retrieve vrack details: %w", err)
	}
	d.Set("plan", serviceObj.ToSDKv2PlanValue())

	// Retrieve subsidiary information
	var me MeResponse
	if err := config.OVHClient.Get("/me", &me); err != nil {
		return fmt.Errorf("error retrieving account information: %w", err)
	}
	d.Set("ovh_subsidiary", me.OvhSubsidiary)

	return nil
}

func resourceVrackDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Id()

	terminate := func() (string, error) {
		log.Printf("[DEBUG] Will terminate vrack %s", serviceName)
		endpoint := fmt.Sprintf(
			"/vrack/%s/terminate",
			url.PathEscape(serviceName),
		)
		if err := config.OVHClient.Post(endpoint, nil, nil); err != nil {
			return "", fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
		}
		return serviceName, nil
	}

	confirmTerminate := func(token string) error {
		log.Printf("[DEBUG] Will confirm termination of vrack %s", serviceName)
		endpoint := fmt.Sprintf(
			"/vrack/%s/confirmTermination",
			url.PathEscape(serviceName),
		)
		if err := config.OVHClient.Post(endpoint, &ConfirmTerminationOpts{Token: token}, nil); err != nil {
			return fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
		}
		return nil
	}

	if err := orderDeleteFromResource(d, meta, terminate, confirmTerminate); err != nil {
		return err
	}

	d.SetId("")

	return nil
}

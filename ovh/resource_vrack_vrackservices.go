package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVrackVrackServices() *schema.Resource {
	return &schema.Resource{
		Create: resourceVrackVrackServicesCreate,
		Read:   resourceVrackVrackServicesRead,
		Delete: resourceVrackVrackServicesDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVrackVrackServicesImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of your vrack",
			},
			"vrack_services": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "vrackServices service name",
			},
		},
	}
}

func resourceVrackVrackServicesImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.Split(givenId, "/")
	if len(splitId) != 2 {
		return nil, fmt.Errorf("import ID is not serviceName/vrackServicesName formatted")
	}
	serviceName := splitId[0]
	vrackServices := splitId[1]

	d.SetId(fmt.Sprintf("vrack_%s-vrackServices_%s", serviceName, vrackServices))
	d.Set("service_name", serviceName)
	d.Set("vrack_services", vrackServices)

	return []*schema.ResourceData{d}, nil
}

func resourceVrackVrackServicesCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	vrackServices := d.Get("vrack_services").(string)

	opts := &VrackVrackServicesCreateOpts{
		VrackServices: vrackServices,
	}
	task := VrackTask{}

	endpoint := fmt.Sprintf("/vrack/%s/vrackServices", url.PathEscape(serviceName))
	if err := config.OVHClient.Post(endpoint, opts, &task); err != nil {
		return fmt.Errorf("error calling POST %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	if err := waitForVrackTask(&task, config.OVHClient); err != nil {
		return fmt.Errorf("error waiting for vrack (%s) to attach vrackServices %v: %s", serviceName, opts, err)
	}

	d.SetId(fmt.Sprintf("vrack_%s-vrackServices_%s", serviceName, vrackServices))

	return resourceVrackVrackServicesRead(d, meta)
}

func resourceVrackVrackServicesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	vrackServices := d.Get("vrack_services").(string)

	endpoint := fmt.Sprintf("/vrack/%s/vrackServices/%s",
		url.PathEscape(serviceName),
		url.PathEscape(vrackServices),
	)

	if err := config.OVHClient.Get(endpoint, nil); err != nil {
		return fmt.Errorf("failed to get vrack-vrackServices link: %w", err)
	}

	return nil
}

func resourceVrackVrackServicesDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	vrackServices := d.Get("vrack_services").(string)
	task := VrackTask{}

	endpoint := fmt.Sprintf("/vrack/%s/vrackServices/%s",
		url.PathEscape(serviceName),
		url.PathEscape(vrackServices),
	)

	if err := config.OVHClient.Delete(endpoint, &task); err != nil {
		return fmt.Errorf("error calling DELETE %s with %s/%s:\n\t %q", endpoint, serviceName, vrackServices, err)
	}

	if err := waitForVrackTask(&task, config.OVHClient); err != nil {
		return fmt.Errorf("error waiting for vrack (%s) to detach vrackServices (%s): %s", serviceName, vrackServices, err)
	}

	d.SetId("")

	return nil
}

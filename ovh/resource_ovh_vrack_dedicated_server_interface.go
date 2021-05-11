package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceVrackDedicatedServerInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceVrackDedicatedServerInterfaceCreate,
		Read:   resourceVrackDedicatedServerInterfaceRead,
		Delete: resourceVrackDedicatedServerInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVrackDedicatedServerInterfaceImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_VRACK_SERVICE", nil),
				Description: "Service name of the resource representing the id of the cloud project.",
			},
			"interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVrackDedicatedServerInterfaceImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not SERVICE_NAME/INTERFACE_ID formatted")
	}
	serviceName := splitId[0]
	interfaceId := splitId[1]
	d.SetId(fmt.Sprintf("vrack_%s-dedicatedserver_%s", serviceName, interfaceId))
	d.Set("service_name", serviceName)
	d.Set("interface_id", interfaceId)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceVrackDedicatedServerInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	opts := (&VrackDedicatedServerInterfaceCreateOpts{}).FromResource(d)
	task := &VrackTask{}

	endpoint := fmt.Sprintf("/vrack/%s/dedicatedServerInterface", serviceName)

	if err := config.OVHClient.Post(endpoint, opts, task); err != nil {
		return fmt.Errorf("Error calling POST %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	if err := waitForVrackTask(task, config.OVHClient); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to attach dedicated server interface %v: %s", serviceName, opts, err)
	}

	//set id
	d.SetId(fmt.Sprintf("vrack_%s-dedicatedserverinterface_%s", serviceName, opts.DedicatedServerInterface))

	return resourceVrackDedicatedServerInterfaceRead(d, meta)
}

func resourceVrackDedicatedServerInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	vds := &VrackDedicatedServerInterface{}

	serviceName := d.Get("service_name").(string)
	interfaceId := d.Get("interface_id").(string)

	endpoint := fmt.Sprintf("/vrack/%s/dedicatedServerInterface/%s",
		url.PathEscape(serviceName),
		url.PathEscape(interfaceId),
	)

	if err := config.OVHClient.Get(endpoint, vds); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("service_name", vds.Vrack)
	d.Set("interface_id", vds.DedicatedServerInterface)
	return nil
}

func resourceVrackDedicatedServerInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	interfaceId := d.Get("interface_id").(string)

	task := &VrackTask{}
	endpoint := fmt.Sprintf("/vrack/%s/dedicatedServerInterface/%s",
		url.PathEscape(serviceName),
		url.PathEscape(interfaceId),
	)

	if err := config.OVHClient.Delete(endpoint, task); err != nil {
		return fmt.Errorf("Error calling DELETE %s with %s/%s:\n\t %q", endpoint, serviceName, interfaceId, err)
	}

	if err := waitForVrackTask(task, config.OVHClient); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to detach dedicated server (%s): %s", serviceName, interfaceId, err)
	}

	d.SetId("")
	return nil
}

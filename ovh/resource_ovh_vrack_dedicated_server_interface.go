package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"vrack_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
		return nil, fmt.Errorf("Import Id is not VRACK_ID/INTERFACE_ID formatted")
	}
	vrackId := splitId[0]
	interfaceId := splitId[1]
	d.SetId(fmt.Sprintf("vrack_%s-dedicatedserver_%s", vrackId, interfaceId))
	d.Set("vrack_id", vrackId)
	d.Set("interface_id", interfaceId)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceVrackDedicatedServerInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	vrackId := d.Get("vrack_id").(string)
	opts := (&VrackDedicatedServerInterfaceCreateOpts{}).FromResource(d)
	task := &VrackTask{}

	endpoint := fmt.Sprintf("/vrack/%s/dedicatedServerInterface", vrackId)

	if err := config.OVHClient.Post(endpoint, opts, task); err != nil {
		return fmt.Errorf("Error calling POST %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	if err := waitForVrackTask(task, config.OVHClient); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to attach dedicated server interface %v: %s", vrackId, opts, err)
	}

	//set id
	d.SetId(fmt.Sprintf("vrack_%s-dedicatedserverinterface_%s", vrackId, opts.DedicatedServerInterface))

	return resourceVrackDedicatedServerInterfaceRead(d, meta)
}

func resourceVrackDedicatedServerInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	vds := &VrackDedicatedServerInterface{}

	vrackId := d.Get("vrack_id").(string)
	interfaceId := d.Get("interface_id").(string)

	endpoint := fmt.Sprintf("/vrack/%s/dedicatedServerInterface/%s",
		url.PathEscape(vrackId),
		url.PathEscape(interfaceId),
	)

	err := config.OVHClient.Get(endpoint, vds)
	if err != nil {
		return err
	}

	d.Set("vrack_id", vds.Vrack)
	d.Set("interface_id", vds.DedicatedServerInterface)

	return nil
}

func resourceVrackDedicatedServerInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	vrackId := d.Get("vrack_id").(string)
	interfaceId := d.Get("interface_id").(string)

	task := &VrackTask{}
	endpoint := fmt.Sprintf("/vrack/%s/dedicatedServerInterface/%s",
		url.PathEscape(vrackId),
		url.PathEscape(interfaceId),
	)

	if err := config.OVHClient.Delete(endpoint, task); err != nil {
		return fmt.Errorf("Error calling DELETE %s with %s/%s:\n\t %q", endpoint, vrackId, interfaceId, err)
	}

	if err := waitForVrackTask(task, config.OVHClient); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to detach dedicated server (%s): %s", vrackId, interfaceId, err)
	}

	d.SetId("")
	return nil
}

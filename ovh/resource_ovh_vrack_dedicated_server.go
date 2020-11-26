package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVrackDedicatedServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceVrackDedicatedServerCreate,
		Read:   resourceVrackDedicatedServerRead,
		Delete: resourceVrackDedicatedServerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVrackDedicatedServerImportState,
		},

		Schema: map[string]*schema.Schema{
			"vrack_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVrackDedicatedServerImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not VRACK_ID/SERVER_ID formatted")
	}
	vrackId := splitId[0]
	serverId := splitId[1]
	d.SetId(fmt.Sprintf("vrack_%s-dedicatedserver_%s", vrackId, serverId))
	d.Set("vrack_id", vrackId)
	d.Set("server_id", serverId)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceVrackDedicatedServerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	vrackId := d.Get("vrack_id").(string)
	opts := (&VrackDedicatedServerCreateOpts{}).FromResource(d)
	task := &VrackTask{}

	endpoint := fmt.Sprintf("/vrack/%s/dedicatedServer", vrackId)

	if err := config.OVHClient.Post(endpoint, opts, task); err != nil {
		return fmt.Errorf("Error calling POST %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	if err := waitForVrackTask(task, config.OVHClient); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to attach dedicated server %v: %s", vrackId, opts, err)
	}

	//set id
	d.SetId(fmt.Sprintf("vrack_%s-dedicatedserver_%s", vrackId, opts.DedicatedServer))

	return resourceVrackDedicatedServerRead(d, meta)
}

func resourceVrackDedicatedServerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	vds := &VrackDedicatedServer{}

	vrackId := d.Get("vrack_id").(string)
	serverId := d.Get("server_id").(string)

	endpoint := fmt.Sprintf("/vrack/%s/dedicatedServer/%s",
		url.PathEscape(vrackId),
		url.PathEscape(serverId),
	)

	err := config.OVHClient.Get(endpoint, vds)
	if err != nil {
		return err
	}

	d.Set("vrack_id", vds.Vrack)
	d.Set("server_id", vds.DedicatedServer)

	return nil
}

func resourceVrackDedicatedServerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	vrackId := d.Get("vrack_id").(string)
	serverId := d.Get("server_id").(string)

	task := &VrackTask{}
	endpoint := fmt.Sprintf("/vrack/%s/dedicatedServer/%s",
		url.PathEscape(vrackId),
		url.PathEscape(serverId),
	)

	if err := config.OVHClient.Delete(endpoint, task); err != nil {
		return fmt.Errorf("Error calling DELETE %s with %s/%s:\n\t %q", endpoint, vrackId, serverId, err)
	}

	if err := waitForVrackTask(task, config.OVHClient); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to detach dedicated server (%s): %s", vrackId, serverId, err)
	}

	d.SetId("")
	return nil
}

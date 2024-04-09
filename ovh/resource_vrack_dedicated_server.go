package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
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
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_VRACK_SERVICE", nil),
				Description: "Service name of the resource representing the id of the cloud project.",
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
		return nil, fmt.Errorf("Import Id is not SERVICE_NAME/SERVER_ID formatted")
	}
	serviceName := splitId[0]
	serverId := splitId[1]
	d.SetId(fmt.Sprintf("vrack_%s-dedicatedserver_%s", serviceName, serverId))
	d.Set("service_name", serviceName)
	d.Set("server_id", serverId)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceVrackDedicatedServerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)

	opts := (&VrackDedicatedServerCreateOpts{}).FromResource(d)
	task := &VrackTask{}

	endpoint := fmt.Sprintf("/vrack/%s/dedicatedServer", serviceName)

	if err := config.OVHClient.Post(endpoint, opts, task); err != nil {
		return fmt.Errorf("Error calling POST %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	if err := waitForVrackTask(task, config.OVHClient); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to attach dedicated server %v: %s", serviceName, opts, err)
	}

	//set id
	d.SetId(fmt.Sprintf("vrack_%s-dedicatedserver_%s", serviceName, opts.DedicatedServer))

	return resourceVrackDedicatedServerRead(d, meta)
}

func resourceVrackDedicatedServerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	vds := &VrackDedicatedServer{}
	serviceName := d.Get("service_name").(string)
	serverId := d.Get("server_id").(string)

	endpoint := fmt.Sprintf("/vrack/%s/dedicatedServer/%s",
		url.PathEscape(serviceName),
		url.PathEscape(serverId),
	)

	if err := config.OVHClient.Get(endpoint, vds); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("service_name", vds.Vrack)
	d.Set("server_id", vds.DedicatedServer)

	return nil
}

func resourceVrackDedicatedServerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	serverId := d.Get("server_id").(string)

	task := &VrackTask{}
	endpoint := fmt.Sprintf("/vrack/%s/dedicatedServer/%s",
		url.PathEscape(serviceName),
		url.PathEscape(serverId),
	)

	if err := config.OVHClient.Delete(endpoint, task); err != nil {
		return fmt.Errorf("Error calling DELETE %s with %s/%s:\n\t %q", endpoint, serviceName, serverId, err)
	}

	if err := waitForVrackTask(task, config.OVHClient); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to detach dedicated server (%s): %s", serviceName, serverId, err)
	}

	d.SetId("")
	return nil
}

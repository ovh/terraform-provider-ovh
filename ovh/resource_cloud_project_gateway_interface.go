package ovh

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceOvhCloudProjectGatewayInterfaceImportState(
	d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not OVH_CLOUD_PROJECT/network_id formatted")
	}
	d.SetId(splitId[1])
	d.Set("service_name", splitId[0])
	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectGatewayInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectGatewayInterfaceCreate,
		Read:   resourceCloudProjectGatewayInterfaceRead,
		Delete: resourceCloudProjectGatewayInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceOvhCloudProjectGatewayInterfaceImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
				Description: "Service name of the resource representing the id of the cloud project.",
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudProjectGatewayInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	region := d.Get("region").(string)
	gatewayId := d.Get("gateway_id").(string)

	params := &CloudProjectGatewayInterfaceCreateOpts{
		Subnet: d.Get("subnet_id").(string),
	}

	r := &CloudProjectGatewayInterfaceResponse{}
	log.Printf("[DEBUG] Will create public cloud gateway interface: %s", params)

	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/gateway/%s/interface", serviceName, region, gatewayId)

	if err := config.OVHClient.Post(endpoint, params, r); err != nil {
		return fmt.Errorf("calling %s with params %s:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Created Gateway interface %s", r)

	//set id
	d.SetId(r.Id)

	return resourceCloudProjectGatewayInterfaceRead(d, meta)
}

func resourceCloudProjectGatewayInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	region := d.Get("region").(string)
	gatewayId := d.Get("gateway_id").(string)

	r := &CloudProjectGatewayInterfaceResponse{}

	log.Printf("[DEBUG] Will read public cloud gateway interface for project: %s, region: %s, gateway: %s, id: %s", serviceName, region, gatewayId, d.Id())

	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/gateway/%s/interface/%s", serviceName, region, gatewayId, d.Id())

	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId(r.Id)
	d.Set("subnet_id", r.SubnetId)
	d.Set("network_id", r.NetworkId)
	d.Set("ip", r.Ip)

	log.Printf("[DEBUG] Read Public Cloud Gateway Interface %s", r)
	return nil
}

func resourceCloudProjectGatewayInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	region := d.Get("region").(string)
	gatewayId := d.Get("gateway_id").(string)

	id := d.Id()

	log.Printf("[DEBUG] Will delete public cloud gateway interface for project: %s, region: %s, id: %s", serviceName, region, id)

	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/gateway/%s/interface/%s", serviceName, region, gatewayId, id)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return fmt.Errorf("calling %s:\n\t %q", endpoint, err)
	}

	d.SetId("")

	log.Printf("[DEBUG] Deleted Public Cloud %s Gateway %s Interface %s", serviceName, gatewayId, id)
	return nil
}

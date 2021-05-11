package ovh

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceOvhCloudProjectNetworkPrivateSubnetImportState(
	d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 3)
	if len(splitId) != 3 {
		return nil, fmt.Errorf("Import Id is not service_name/network_id/subnet_id formatted")
	}
	d.SetId(splitId[2])
	d.Set("network_id", splitId[1])
	d.Set("service_name", splitId[0])
	results := make([]*schema.ResourceData, 1)
	results[0] = d
	log.Printf(
		"[DEBUG] Will Import ovh_cloud_project_network_private_subnet with project %s, network %s, id %s",
		d.Get("service_name"),
		d.Get("network_id"),
		d.Id(),
	)
	return results, nil
}

func resourceCloudProjectNetworkPrivateSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectNetworkPrivateSubnetCreate,
		Read:   resourceCloudProjectNetworkPrivateSubnetRead,
		Delete: resourceCloudProjectNetworkPrivateSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: resourceOvhCloudProjectNetworkPrivateSubnetImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
				Description: "Service name of the resource representing the id of the cloud project.",
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dhcp": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"start": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: resourceCloudProjectNetworkPrivateSubnetValidateIP,
			},
			"end": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: resourceCloudProjectNetworkPrivateSubnetValidateIP,
			},
			"network": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: resourceCloudProjectNetworkPrivateSubnetValidateNetwork,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"no_gateway": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"gateway_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_pools": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dhcp": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"end": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceCloudProjectNetworkPrivateSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	networkId := d.Get("network_id").(string)

	params := &CloudProjectNetworkPrivatesCreateOpts{
		ServiceName: serviceName,
		NetworkId:   networkId,
		Dhcp:        d.Get("dhcp").(bool),
		NoGateway:   d.Get("no_gateway").(bool),
		Start:       d.Get("start").(string),
		End:         d.Get("end").(string),
		Network:     d.Get("network").(string),
		Region:      d.Get("region").(string),
	}

	r := &CloudProjectNetworkPrivatesResponse{}

	log.Printf("[DEBUG] Will create public cloud private network subnet: %s", params)

	endpoint := fmt.Sprintf("/cloud/project/%s/network/private/%s/subnet", serviceName, networkId)

	if err := config.OVHClient.Post(endpoint, params, r); err != nil {
		return fmt.Errorf("calling POST %s with params %s:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Created Private Network Subnet %s", r)

	//set id
	d.SetId(r.Id)

	return resourceCloudProjectNetworkPrivateSubnetRead(d, meta)
}

func resourceCloudProjectNetworkPrivateSubnetRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	networkId := d.Get("network_id").(string)

	subnets := []*CloudProjectNetworkPrivatesResponse{}

	log.Printf("[DEBUG] Will read public cloud private network subnet for project: %s, network: %s, id: %s", serviceName, networkId, d.Id())

	endpoint := fmt.Sprintf("/cloud/project/%s/network/private/%s/subnet", serviceName, networkId)

	if err := config.OVHClient.Get(endpoint, &subnets); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	var r *CloudProjectNetworkPrivatesResponse
	for i := range subnets {
		if subnets[i].Id == d.Id() {
			r = subnets[i]
		}
	}

	if r == nil {
		d.SetId("")
		return nil
	}

	d.Set("gateway_ip", r.GatewayIp)
	d.Set("cidr", r.Cidr)

	ippools := make([]map[string]interface{}, 0)
	for i := range r.IPPools {
		ippool := make(map[string]interface{})
		ippool["network"] = r.IPPools[i].Network
		ippool["region"] = r.IPPools[i].Region
		ippool["dhcp"] = r.IPPools[i].Dhcp
		ippool["start"] = r.IPPools[i].Start
		ippool["end"] = r.IPPools[i].End
		ippools = append(ippools, ippool)
	}

	d.Set("network", ippools[0]["network"])
	d.Set("region", ippools[0]["region"])
	d.Set("dhcp", ippools[0]["dhcp"])
	d.Set("start", ippools[0]["start"])
	d.Set("end", ippools[0]["end"])
	d.Set("ip_pools", ippools)

	if r.GatewayIp == "" {
		d.Set("no_gateway", true)
	} else {
		d.Set("no_gateway", false)
	}

	d.SetId(r.Id)

	d.Set("service_name", serviceName)
	log.Printf("[DEBUG] Read Public Cloud Private Network %v", r)
	return nil
}

func resourceCloudProjectNetworkPrivateSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	networkId := d.Get("network_id").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will delete public cloud private network subnet for project: %s, network: %s, id: %s", serviceName, networkId, id)

	endpoint := fmt.Sprintf("/cloud/project/%s/network/private/%s/subnet/%s", serviceName, id, id)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return fmt.Errorf("calling DELETE %s:\n\t %q", endpoint, err)
	}

	d.SetId("")

	log.Printf("[DEBUG] Deleted Public Cloud %s Private Network %s Subnet %s", serviceName, networkId, id)
	return nil
}

func resourceCloudProjectNetworkPrivateSubnetValidateIP(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	ip := net.ParseIP(value)
	if ip == nil {
		errors = append(errors, fmt.Errorf("%q must be a valid IP", k))
	}
	return
}

func resourceCloudProjectNetworkPrivateSubnetValidateNetwork(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, _, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q is not a valid network value: %#v", k, err))
	}
	return
}

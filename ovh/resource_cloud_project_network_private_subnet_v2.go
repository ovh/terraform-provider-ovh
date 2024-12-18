package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceOvhCloudProjectNetworkPrivateSubnetV2ImportState(
	d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 4)
	if len(splitId) != 4 {
		return nil, fmt.Errorf("Import Id is not service_name/region/network_id/subnet_id formatted")
	}
	d.Set("service_name", splitId[0])
	d.Set("region", splitId[1])
	d.Set("network_id", splitId[2])
	d.SetId(splitId[3])
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

func resourceCloudProjectNetworkPrivateSubnetV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectNetworkPrivateSubnetV2Create,
		Read:   resourceCloudProjectNetworkPrivateSubnetV2Read,
		Delete: resourceCloudProjectNetworkPrivateSubnetV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceOvhCloudProjectNetworkPrivateSubnetV2ImportState,
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Network ID of subnet",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region of network/subnet",
			},
			"dhcp": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     true,
				Description: "Enable DHCP in subnet",
			},
			"cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: resourceCloudProjectNetworkPrivateSubnetValidateNetwork,
				Description:  "CIDR of subnet",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of subnet",
			},
			"gateway_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				Description:  "Gateway IP of subnet",
				ValidateFunc: resourceCloudProjectNetworkPrivateSubnetValidateIP,
			},
			"enable_gateway_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     true,
				Description: "Enable gateway IP in subnet",
			},
			"use_default_public_dns_resolver": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Use OVH default DNS",
			},
			"dns_nameservers": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: resourceCloudProjectNetworkPrivateSubnetValidateIP,
				},
				Description: "List of DNS nameservers, default: 213.186.33.99",
			},
			"host_route": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: resourceCloudProjectNetworkPrivateSubnetValidateNetwork,
						},
						"nexthop": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: resourceCloudProjectNetworkPrivateSubnetValidateIP,
						},
					},
				},
				Description: "Static host routes of subnet",
			},
			"allocation_pool": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: resourceCloudProjectNetworkPrivateSubnetValidateIP,
						},
						"end": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: resourceCloudProjectNetworkPrivateSubnetValidateIP,
						},
					},
				},
				Description: "DHCP allocation pools of subnet",
			},
		},
	}
}

func resourceCloudProjectNetworkPrivateSubnetV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	regionName := d.Get("region").(string)
	networkId := d.Get("network_id").(string)
	subnetName := d.Get("name").(string)
	cidr := d.Get("cidr").(string)
	enableGatewayIP := d.Get("enable_gateway_ip").(bool)
	enableDHCP := d.Get("dhcp").(bool)
	useDefaultPublicDNSResolverbool, ok := d.GetOkExists("use_default_public_dns_resolver")
	var useDefaultPublicDNSResolver *bool
	if ok {
		i := useDefaultPublicDNSResolverbool.(bool)
		useDefaultPublicDNSResolver = &i
	}
	gatewayIp := d.Get("gateway_ip").(string)

	hostRoutesStrings, err := helpers.StringMapFromSchema(d, "host_route", "destination", "nexthop")
	if err != nil {
		return err
	}

	hostRoutes := []HostRoute{}
	for _, hostRouteStrings := range hostRoutesStrings {
		hostRoutes = append(hostRoutes, HostRoute{
			Destination: hostRouteStrings["destination"],
			Nexthop:     hostRouteStrings["nexthop"],
		})
	}

	allocationPoolsStrings, err := helpers.StringMapFromSchema(d, "allocation_pool", "start", "end")
	if err != nil {
		return err
	}

	allocationPools := []AllocationPool{}
	for _, allocationPoolStrings := range allocationPoolsStrings {
		allocationPools = append(allocationPools, AllocationPool{
			Start: allocationPoolStrings["start"],
			End:   allocationPoolStrings["end"],
		})
	}

	dnsNameServers, err := helpers.StringsFromSchema(d, "dns_nameservers")
	if err != nil {
		return err
	}

	createSubnetParams := &CloudProjectNetworkPrivateV2CreateOpts{
		Name:                        subnetName,
		Cidr:                        cidr,
		IpVersion:                   4,
		AllocationPools:             allocationPools,
		DnsNameServers:              dnsNameServers,
		GatewayIp:                   gatewayIp,
		EnableGatewayIP:             enableGatewayIP,
		EnableDHCP:                  enableDHCP,
		HostRoutes:                  hostRoutes,
		UseDefaultPublicDNSResolver: useDefaultPublicDNSResolver,
	}

	subnetResponse := &CloudProjectNetworkPrivateV2Response{}
	endpointPostSubnet := fmt.Sprintf("/cloud/project/%s/region/%s/network/%s/subnet",
		url.PathEscape(serviceName),
		url.PathEscape(regionName),
		url.PathEscape(networkId),
	)
	err = config.OVHClient.Post(endpointPostSubnet, createSubnetParams, subnetResponse)
	if err != nil {
		return fmt.Errorf("calling POST %s with params %v:\n\t %q", endpointPostSubnet, createSubnetParams, err)
	}

	d.Set("gateway_ip", subnetResponse.GatewayIp)
	d.Set("cidr", subnetResponse.Cidr)
	d.Set("enable_gateway_ip", subnetResponse.GatewayIp != nil)
	d.Set("service_name", serviceName)
	d.Set("network_id", networkId)
	d.Set("dhcp", subnetResponse.DHCPEnabled)
	d.Set("region", regionName)
	d.SetId(subnetResponse.Id)
	log.Printf("[DEBUG] Read Public Cloud Private Network %v", subnetResponse)
	return nil
}

func resourceCloudProjectNetworkPrivateSubnetV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	networkId := d.Get("network_id").(string)
	regionName := d.Get("region").(string)

	log.Printf("[DEBUG] Will read public cloud private network subnet for project: %s, region: %s, network: %s, id: %s", serviceName, regionName, networkId, d.Id())

	subnet := CloudProjectNetworkPrivateV2Response{}
	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/network/%s/subnet/%s", serviceName, regionName, networkId, d.Id())
	if err := config.OVHClient.Get(endpoint, &subnet); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("gateway_ip", subnet.GatewayIp)
	d.Set("cidr", subnet.Cidr)
	d.Set("dhcp", subnet.DHCPEnabled)
	d.Set("enable_gateway_ip", subnet.GatewayIp != nil)
	d.Set("service_name", serviceName)
	d.Set("network_id", networkId)
	d.Set("region", regionName)
	d.SetId(subnet.Id)
	log.Printf("[DEBUG] Read Public Cloud Private Network %v", subnet)
	return nil
}

func resourceCloudProjectNetworkPrivateSubnetV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	networkId := d.Get("network_id").(string)
	regionName := d.Get("region").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will delete public cloud private network subnet V2 for project: %s, region: %s, network: %s, id: %s", serviceName, regionName, networkId, id)

	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/network/%s/subnet/%s",
		url.PathEscape(serviceName),
		url.PathEscape(regionName),
		url.PathEscape(networkId),
		url.PathEscape(id),
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return fmt.Errorf("calling DELETE %s:\n\t %q", endpoint, err)
	}

	d.SetId("")

	log.Printf("[DEBUG] Deleted Public Cloud %s Private Network %s Subnet %s", serviceName, networkId, id)
	return nil
}

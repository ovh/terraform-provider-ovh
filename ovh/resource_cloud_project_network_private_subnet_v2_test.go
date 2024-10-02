package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccCloudProjectNetworkPrivateV2SubnetConfig(config string) string {
	var testAccCloudProjectNetworkPrivateSubnetV2Config_attachVrack = `
	resource "ovh_vrack_cloudproject" "attach" {
		service_name = "%s"
		project_id   = "%s"
	}

	data "ovh_cloud_project_regions" "regions" {
		service_name = ovh_vrack_cloudproject.attach.project_id
		has_services_up = ["network"]
	}

	resource "ovh_cloud_project_network_private" "network" {
		service_name = ovh_vrack_cloudproject.attach.project_id
		vlan_id		= 1
		name		= "terraform_testacc_private_net"
		regions		= slice(sort(tolist(data.ovh_cloud_project_regions.regions.names)), 0, 3)
	}
`
	attachVrack := fmt.Sprintf(
		testAccCloudProjectNetworkPrivateSubnetV2Config_attachVrack,
		os.Getenv("OVH_VRACK_SERVICE_TEST"),
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
	)

	var testAccCloudProjectNetworkPrivateSubnetV2Config_noAttachVrack = `
	data "ovh_cloud_project_regions" "regions" {
		service_name = "%s"
		has_services_up = ["network"]
	}

	resource "ovh_cloud_project_network_private" "network" {
		service_name = data.ovh_cloud_project_regions.regions.service_name
		vlan_id	= 1
		name	= "terraform_testacc_private_net"
		regions	= slice(sort(tolist(data.ovh_cloud_project_regions.regions.names)), 0, 3)
	}
`
	noAttachVrack := fmt.Sprintf(
		testAccCloudProjectNetworkPrivateSubnetV2Config_noAttachVrack,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
	)

	if os.Getenv("OVH_ATTACH_VRACK") == "0" {
		return fmt.Sprintf(
			config,
			noAttachVrack,
		)
	}

	return fmt.Sprintf(
		config,
		attachVrack,
	)
}

var testAccCloudProjectNetworkPrivateSubnetV2Config_basic = `
%s

resource "ovh_cloud_project_network_private_subnet_v2" "subnet" {
  service_name = ovh_cloud_project_network_private.network.service_name

  # whatever region, for test purpose
  network_id        = element(tolist(ovh_cloud_project_network_private.network.regions_attributes), 0).openstackid
  region            = element(tolist(ovh_cloud_project_network_private.network.regions), 0)
  name              = "my_new_subnet"
  cidr              = "192.168.169.0/24"
  dns_nameservers   = ["1.1.1.1"]
  host_route {
    destination = "192.168.169.0/24" 
    nexthop = "192.168.169.254"
  }
  allocation_pool {
    start = "192.168.169.100"
    end = "192.168.169.200"
  }
  dhcp              = true
  gateway_ip        = "192.168.169.253"
  enable_gateway_ip = true
}
`

func TestAccCloudProjectNetworkPrivateSubnetV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckVRack(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudProjectNetworkPrivateV2SubnetConfig(testAccCloudProjectNetworkPrivateSubnetV2Config_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_network_private_subnet_v2.subnet", "service_name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_network_private_subnet_v2.subnet", "network_id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_network_private_subnet_v2.subnet", "region"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "gateway_ip", "192.168.169.253"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "allocation_pool.0.start", "192.168.169.100"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "allocation_pool.0.end", "192.168.169.200"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "host_route.0.destination", "192.168.169.0/24"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "host_route.0.nexthop", "192.168.169.254"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "dns_nameservers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "dhcp", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "enable_gateway_ip", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "cidr", "192.168.169.0/24"),
				),
			},
		},
	})
}

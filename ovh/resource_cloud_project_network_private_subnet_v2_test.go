package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectNetworkPrivateSubnetV2_basic(t *testing.T) {

	cloudProject := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	var testAccCloudProjectNetworkPrivateSubnetV2Config_basic = fmt.Sprintf(`
resource "ovh_cloud_project_network_private" "network" {
  service_name = "%s"
  vlan_id    = 0
  name       = "terraform_testacc_private_net"
  regions    = ["%s"]
}

  resource "ovh_cloud_project_network_private_subnet_v2" "subnet" {
  service_name = ovh_cloud_project_network_private.network.service_name

  network_id        = element(tolist(ovh_cloud_project_network_private.network.regions_attributes), 0).openstackid
  region     = element(tolist(sort(ovh_cloud_project_network_private.network.regions)), 0)
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
  use_default_public_dns_resolver = false
}
`, cloudProject, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckVRack(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudProjectNetworkPrivateSubnetV2Config_basic,
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
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "use_default_public_dns_resolver", "false"),
				),
			},
		},
	})
}

func TestAccCloudProjectNetworkPrivateSubnetV2DefaultDns_basic(t *testing.T) {

	cloudProject := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	var testAccCloudProjectNetworkPrivateSubnetV2Config_basic = fmt.Sprintf(`
resource "ovh_cloud_project_network_private" "network" {
  service_name = "%s"
  vlan_id    = 0
  name       = "terraform_testacc_private_net"
  regions    = ["%s"]
}

  resource "ovh_cloud_project_network_private_subnet_v2" "subnet" {
  service_name = ovh_cloud_project_network_private.network.service_name

  network_id        = element(tolist(ovh_cloud_project_network_private.network.regions_attributes), 0).openstackid
  region     = element(tolist(sort(ovh_cloud_project_network_private.network.regions)), 0)
  name              = "my_new_subnet"
  cidr              = "192.168.169.0/24"
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
  use_default_public_dns_resolver = true
}
`, cloudProject, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckVRack(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudProjectNetworkPrivateSubnetV2Config_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_network_private_subnet_v2.subnet", "service_name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_network_private_subnet_v2.subnet", "network_id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_network_private_subnet_v2.subnet", "region"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "gateway_ip", "192.168.169.253"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "allocation_pool.0.start", "192.168.169.100"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "allocation_pool.0.end", "192.168.169.200"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "host_route.0.destination", "192.168.169.0/24"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "host_route.0.nexthop", "192.168.169.254"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "dns_nameservers.0", "213.186.33.99"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "dhcp", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "enable_gateway_ip", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "cidr", "192.168.169.0/24"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "use_default_public_dns_resolver", "true"),
				),
			},
		},
	})
}

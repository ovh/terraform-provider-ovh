package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testAccPublicCloudPrivateNetworkSubnetConfig_attachVrack = `
resource "ovh_vrack_cloudproject" "attach" {
  vrack_id   = "%s"
  project_id = "%s"
}

data "ovh_cloud_regions" "regions" {
  project_id = ovh_vrack_cloudproject.attach.project_id

  has_services_up = ["network"]
}
`

var testAccPublicCloudPrivateNetworkSubnetConfig_noAttachVrack = `
data "ovh_cloud_regions" "regions" {
  project_id = "%s"

  has_services_up = ["network"]
}
`

var testAccPublicCloudPrivateNetworkSubnetConfig_basic = `
%s

resource "ovh_cloud_network_private" "network" {
  project_id = data.ovh_cloud_regions.regions.project_id
  vlan_id    = 0
  name       = "terraform_testacc_private_net"
  regions    = tolist(data.ovh_cloud_regions.regions.names)
}

resource "ovh_cloud_network_private_subnet" "subnet" {
  project_id = ovh_cloud_network_private.network.project_id
  network_id = ovh_cloud_network_private.network.id

  # whatever region, for test purpose
  region     = element(tolist(sort(data.ovh_cloud_regions.regions.names)), 0)
  start      = "192.168.168.100"
  end        = "192.168.168.200"
  network    = "192.168.168.0/24"
  dhcp       = true
  no_gateway = false
}
`

func testAccPublicCloudPrivateNetworkSubnetConfig() string {
	attachVrack := fmt.Sprintf(
		testAccPublicCloudPrivateNetworkSubnetConfig_attachVrack,
		os.Getenv("OVH_VRACK"),
		os.Getenv("OVH_PUBLIC_CLOUD"),
	)
	noAttachVrack := fmt.Sprintf(
		testAccPublicCloudPrivateNetworkSubnetConfig_noAttachVrack,
		os.Getenv("OVH_PUBLIC_CLOUD"),
	)

	if os.Getenv("OVH_ATTACH_VRACK") == "0" {
		return fmt.Sprintf(
			testAccPublicCloudPrivateNetworkSubnetConfig_basic,
			noAttachVrack,
		)
	}

	return fmt.Sprintf(
		testAccPublicCloudPrivateNetworkSubnetConfig_basic,
		attachVrack,
	)
}

func TestAccPublicCloudPrivateNetworkSubnet_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckPublicCloudPrivateNetworkSubnetPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPublicCloudPrivateNetworkSubnetConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_subnet.subnet", "project_id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_subnet.subnet", "network_id"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_subnet.subnet", "start", "192.168.168.100"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_subnet.subnet", "end", "192.168.168.200"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_subnet.subnet", "network", "192.168.168.0/24"),
				),
			},
		},
	})
}

func testAccCheckPublicCloudPrivateNetworkSubnetPreCheck(t *testing.T) {
	testAccPreCheckPublicCloud(t)
	testAccCheckPublicCloudExists(t)
	testAccPreCheckVRack(t)
}

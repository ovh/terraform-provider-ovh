package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccPublicCloudPrivateNetworkSubnetConfig = fmt.Sprintf(`
resource "ovh_vrack_cloudproject" "attach" {
  vrack_id   = "%s"
  project_id = "%s"
}

data "ovh_cloud_regions" "regions" {
  project_id = ovh_vrack_cloudproject.attach.project_id
}

resource "ovh_cloud_network_private" "network" {
  project_id = ovh_vrack_cloudproject.attach.project_id
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
`, os.Getenv("OVH_VRACK"), os.Getenv("OVH_PUBLIC_CLOUD"))

func TestAccPublicCloudPrivateNetworkSubnet_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccCheckPublicCloudPrivateNetworkSubnetPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPublicCloudPrivateNetworkSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPublicCloudPrivateNetworkSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVRackPublicCloudAttachmentExists("ovh_vrack_cloudproject.attach", t),
					testAccCheckPublicCloudPrivateNetworkExists("ovh_cloud_network_private.network", t),
					testAccCheckPublicCloudPrivateNetworkSubnetExists("ovh_cloud_network_private_subnet.subnet", t),
				),
			},
		},
	})
}

func testAccCheckPublicCloudPrivateNetworkSubnetPreCheck(t *testing.T) {
	testAccPreCheckPublicCloud(t)
	testAccCheckPublicCloudExists(t)
}

func testAccCheckPublicCloudPrivateNetworkSubnetExists(n string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("No Project ID is set")
		}

		if rs.Primary.Attributes["network_id"] == "" {
			return fmt.Errorf("No Network ID is set")
		}

		return publicCloudPrivateNetworkSubnetExists(
			rs.Primary.Attributes["project_id"],
			rs.Primary.Attributes["network_id"],
			rs.Primary.ID,
			config.OVHClient,
		)
	}
}

func testAccCheckPublicCloudPrivateNetworkSubnetDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_cloud_network_private_subnet" {
			continue
		}

		err := publicCloudPrivateNetworkSubnetExists(
			rs.Primary.Attributes["project_id"],
			rs.Primary.Attributes["network_id"],
			rs.Primary.ID,
			config.OVHClient,
		)

		if err == nil {
			return fmt.Errorf("VRack > Public Cloud Private Network Subnet still exists")
		}

	}
	return nil
}

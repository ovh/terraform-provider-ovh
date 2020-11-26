package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAccCloudNetworkPrivateSubnetConfig_attachVrack = `
resource "ovh_vrack_cloudproject" "attach" {
  vrack_id   = "%s"
  project_id = "%s"
}

data "ovh_cloud_regions" "regions" {
  project_id = ovh_vrack_cloudproject.attach.project_id

  has_services_up = ["network"]
}

resource "ovh_cloud_network_private" "network" {
  project_id = ovh_vrack_cloudproject.attach.project_id
  vlan_id    = 0
  name       = "terraform_testacc_private_net"
  regions    = tolist(data.ovh_cloud_regions.regions.names)
}
`

var testAccCloudNetworkPrivateSubnetConfig_noAttachVrack = `
data "ovh_cloud_regions" "regions" {
  project_id = "%s"

  has_services_up = ["network"]
}

resource "ovh_cloud_network_private" "network" {
  project_id = data.ovh_cloud_regions.regions.project_id
  vlan_id    = 0
  name       = "terraform_testacc_private_net"
  regions    = tolist(data.ovh_cloud_regions.regions.names)
}
`

var testAccCloudNetworkPrivateSubnetConfig_basic = `
%s

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

func testAccCloudNetworkPrivateSubnetConfig() string {
	attachVrack := fmt.Sprintf(
		testAccCloudNetworkPrivateSubnetConfig_attachVrack,
		os.Getenv("OVH_VRACK"),
		os.Getenv("OVH_PUBLIC_CLOUD"),
	)
	noAttachVrack := fmt.Sprintf(
		testAccCloudNetworkPrivateSubnetConfig_noAttachVrack,
		os.Getenv("OVH_PUBLIC_CLOUD"),
	)

	if os.Getenv("OVH_ATTACH_VRACK") == "0" {
		return fmt.Sprintf(
			testAccCloudNetworkPrivateSubnetConfig_basic,
			noAttachVrack,
		)
	}

	return fmt.Sprintf(
		testAccCloudNetworkPrivateSubnetConfig_basic,
		attachVrack,
	)
}

func TestAccCloudNetworkPrivateSubnet_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckcCloudNetworkPrivateSubnetPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudNetworkPrivateSubnetConfig(),
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

func testAccCheckcCloudNetworkPrivateSubnetPreCheck(t *testing.T) {
	testAccPreCheckCloud(t)
	testAccCheckCloudExists(t)
	testAccPreCheckVRack(t)
}

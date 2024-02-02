package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccCloudProjectNetworkPrivateSubnetConfig_attachVrack = `
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
  vlan_id    = 0
  name       = "terraform_testacc_private_net"
  regions    = slice(sort(tolist(data.ovh_cloud_project_regions.regions.names)), 0, 3)
}
`

var testAccCloudProjectNetworkPrivateSubnetConfig_noAttachVrack = `
data "ovh_cloud_project_regions" "regions" {
  service_name = "%s"

  has_services_up = ["network"]
}

resource "ovh_cloud_project_network_private" "network" {
  service_name = data.ovh_cloud_project_regions.regions.service_name
  vlan_id    = 0
  name       = "terraform_testacc_private_net"
  regions    = slice(sort(tolist(data.ovh_cloud_project_regions.regions.names)), 0, 3)
}
`

var testAccCloudProjectNetworkPrivateSubnetConfig_basic = `
%s

resource "ovh_cloud_project_network_private_subnet" "subnet" {
  service_name = ovh_cloud_project_network_private.network.service_name
  network_id = ovh_cloud_project_network_private.network.id

  # whatever region, for test purpose
  region     = element(tolist(sort(data.ovh_cloud_project_regions.regions.names)), 0)
  start      = "192.168.168.100"
  end        = "192.168.168.200"
  network    = "192.168.168.0/24"
  dhcp       = true
  no_gateway = false
}
`

func testAccCloudProjectNetworkPrivateSubnetConfig(config string) string {
	attachVrack := fmt.Sprintf(
		testAccCloudProjectNetworkPrivateSubnetConfig_attachVrack,
		os.Getenv("OVH_VRACK_SERVICE_TEST"),
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
	)
	noAttachVrack := fmt.Sprintf(
		testAccCloudProjectNetworkPrivateSubnetConfig_noAttachVrack,
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

func TestAccCloudProjectNetworkPrivateSubnet_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckcCloudProjectNetworkPrivateSubnetPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudProjectNetworkPrivateSubnetConfig(testAccCloudProjectNetworkPrivateSubnetConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_network_private_subnet.subnet", "service_name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_network_private_subnet.subnet", "network_id"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet.subnet", "start", "192.168.168.100"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet.subnet", "end", "192.168.168.200"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet.subnet", "network", "192.168.168.0/24"),
				),
			},
		},
	})
}

func testAccCheckcCloudProjectNetworkPrivateSubnetPreCheck(t *testing.T) {
	testAccPreCheckCloud(t)
	testAccCheckCloudProjectExists(t)
	testAccPreCheckVRack(t)
}

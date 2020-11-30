package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAccCloudNetworkPrivateSubnetConfig_attachVrack = `
resource "ovh_vrack_cloudproject" "attach" {
  service_name = "%s"
  project_id   = "%s"
}

data "ovh_cloud_regions" "regions" {
  service_name = ovh_vrack_cloudproject.attach.project_id

  has_services_up = ["network"]
}

resource "ovh_cloud_network_private" "network" {
  service_name = ovh_vrack_cloudproject.attach.project_id
  vlan_id    = 0
  name       = "terraform_testacc_private_net"
  regions    = slice(sort(tolist(data.ovh_cloud_regions.regions.names)), 0, 3)
}
`

var testAccCloudNetworkPrivateSubnetConfig_noAttachVrack = `
data "ovh_cloud_regions" "regions" {
  service_name = "%s"

  has_services_up = ["network"]
}

resource "ovh_cloud_network_private" "network" {
  service_name = data.ovh_cloud_regions.regions.service_name
  vlan_id    = 0
  name       = "terraform_testacc_private_net"
  regions    = slice(sort(tolist(data.ovh_cloud_regions.regions.names)), 0, 3)
}
`

var testAccCloudNetworkPrivateSubnetConfig_basic = `
%s

resource "ovh_cloud_network_private_subnet" "subnet" {
  service_name = ovh_cloud_network_private.network.service_name
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

var testAccCloudNetworkPrivateSubnetDeprecatedConfig_basic = `
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

func testAccCloudNetworkPrivateSubnetConfig(config string) string {
	attachVrack := fmt.Sprintf(
		testAccCloudNetworkPrivateSubnetConfig_attachVrack,
		os.Getenv("OVH_VRACK_SERVICE_TEST"),
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
	)
	noAttachVrack := fmt.Sprintf(
		testAccCloudNetworkPrivateSubnetConfig_noAttachVrack,
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

func TestAccCloudNetworkPrivateSubnet_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckcCloudNetworkPrivateSubnetPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudNetworkPrivateSubnetConfig(testAccCloudNetworkPrivateSubnetConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_subnet.subnet", "service_name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_subnet.subnet", "network_id"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_subnet.subnet", "start", "192.168.168.100"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_subnet.subnet", "end", "192.168.168.200"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_subnet.subnet", "network", "192.168.168.0/24"),
				),
			},
		},
	})
}

func TestAccCloudNetworkPrivateSubnetDeprecated_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckcCloudNetworkPrivateSubnetPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudNetworkPrivateSubnetConfig(testAccCloudNetworkPrivateSubnetDeprecatedConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_subnet.subnet", "service_name"),
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

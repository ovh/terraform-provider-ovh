package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceCloudNetworkPrivateSubnetNamePrefix = "tf-test-subnet-"

func testAccPreCheckCloudNetworkPrivateSubnet(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST must be set for acceptance tests")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_REGION_TEST must be set for acceptance tests")
	}
	// Creating a regional private network requires a vRack.
	testAccPreCheckVRack(t)
}

func testAccCloudNetworkPrivateSubnetImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s/%s",
			rs.Primary.Attributes["service_name"],
			rs.Primary.Attributes["network_id"],
			rs.Primary.Attributes["id"],
		), nil
	}
}

func TestAccCloudNetworkPrivateVrackSubnet_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	subnetName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateSubnetNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "terraform_testacc_private_net"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "test" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id = ovh_cloud_network_private_vrack.network.id

  name = "%s"
  cidr = "10.0.0.0/24"
  region       = "%s"
}
`, serviceName, region, subnetName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateSubnet(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "network_id"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "name", subnetName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "cidr", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "current_state.name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "current_state.cidr"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "current_state.location.region"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_network_private_vrack_subnet.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudNetworkPrivateSubnetImportStateIdFunc("ovh_cloud_network_private_vrack_subnet.test"),
			},
		},
	})
}

// TestAccCloudNetworkPrivateVrackSubnet_withAllOptions creates a subnet exercising all
// optional attributes: dhcp_enabled, gateway_ip, dns_nameservers and
// allocation_pools.
func TestAccCloudNetworkPrivateVrackSubnet_withAllOptions(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	subnetName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateSubnetNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "terraform_testacc_private_net"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "test" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id = ovh_cloud_network_private_vrack.network.id

  name            = "%s"
  cidr            = "10.0.0.0/24"
  dhcp_enabled    = true
  gateway_ip      = "10.0.0.1"
  dns_nameservers = ["1.1.1.1", "8.8.8.8"]
  region       = "%s"
  allocation_pools = [
    {
      start = "10.0.0.10"
      end   = "10.0.0.100"
    },
  ]
}
`, serviceName, region, subnetName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateSubnet(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "name", subnetName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "cidr", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dhcp_enabled", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "gateway_ip", "10.0.0.1"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dns_nameservers.#", "2"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dns_nameservers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dns_nameservers.1", "8.8.8.8"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "allocation_pools.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "allocation_pools.0.start", "10.0.0.10"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "allocation_pools.0.end", "10.0.0.100"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "current_state.gateway_ip", "10.0.0.1"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "current_state.dhcp_enabled", "true"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_network_private_vrack_subnet.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudNetworkPrivateSubnetImportStateIdFunc("ovh_cloud_network_private_vrack_subnet.test"),
			},
		},
	})
}

// TestAccCloudNetworkPrivateVrackSubnet_updateMutableFields toggles dhcp_enabled and
// updates dns_nameservers / description / name in place (all mutable, no replacement).
func TestAccCloudNetworkPrivateVrackSubnet_updateMutableFields(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	subnetName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateSubnetNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateSubnetNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "terraform_testacc_private_net"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "test" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id = ovh_cloud_network_private_vrack.network.id

  name            = "%s"
  cidr            = "10.0.0.0/24"
  description     = "initial description"
  dhcp_enabled    = false
  dns_nameservers = ["1.1.1.1"]
  region       = "%s"
}
`, serviceName, region, subnetName, region)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "terraform_testacc_private_net"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "test" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id = ovh_cloud_network_private_vrack.network.id

  name            = "%s"
  cidr            = "10.0.0.0/24"
  description     = "updated description"
  dhcp_enabled    = true
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
  region       = "%s"
}
`, serviceName, region, updatedName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateSubnet(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "name", subnetName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "description", "initial description"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dhcp_enabled", "false"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dns_nameservers.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dns_nameservers.0", "1.1.1.1"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "id"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "name", updatedName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "description", "updated description"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dhcp_enabled", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dns_nameservers.#", "2"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dns_nameservers.0", "8.8.8.8"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dns_nameservers.1", "8.8.4.4"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "current_state.dhcp_enabled", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "resource_status", "READY"),
				),
			},
		},
	})
}

func TestAccCloudNetworkPrivateVrackSubnet_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	subnetName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateSubnetNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateSubnetNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "terraform_testacc_private_net"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "test" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id = ovh_cloud_network_private_vrack.network.id

  name        = "%s"
  cidr        = "10.0.0.0/24"
  description = "initial description"
  region       = "%s"
}
`, serviceName, region, subnetName, region)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "terraform_testacc_private_net"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "test" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id = ovh_cloud_network_private_vrack.network.id

  name        = "%s"
  cidr        = "10.0.0.0/24"
  description = "updated description"
  region       = "%s"
}
`, serviceName, region, updatedName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateSubnet(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "name", subnetName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "description", "initial description"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "checksum"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "name", updatedName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "description", "updated description"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "checksum"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "resource_status", "READY"),
				),
			},
		},
	})
}

package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceCloudNetworkPrivateVrackNamePrefix = "tf-test-net-v2-"

func testAccPreCheckCloudNetworkPrivateVrack(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST must be set for acceptance tests")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_REGION_TEST must be set for acceptance tests")
	}
}

func testAccCloudNetworkPrivateVrackImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}

func TestAccCloudNetworkPrivateVrack_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	networkName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}
`, serviceName, networkName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateVrack(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "name", networkName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "current_state.name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "current_state.location.region"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_network_private_vrack.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudNetworkPrivateVrackImportStateIdFunc("ovh_cloud_network_private_vrack.test"),
			},
		},
	})
}

func TestAccCloudNetworkPrivateVrack_withDescription(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	networkName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  description  = "network created by acceptance test"
}
`, serviceName, networkName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateVrack(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "name", networkName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "description", "network created by acceptance test"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "id"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "current_state.description", "network created by acceptance test"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_network_private_vrack.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudNetworkPrivateVrackImportStateIdFunc("ovh_cloud_network_private_vrack.test"),
			},
		},
	})
}

func TestAccCloudNetworkPrivateVrack_withVlanId(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	networkName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackNamePrefix)
	vlanId := acctest.RandIntRange(0, 4094)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  vlan_id = %d
}
`, serviceName, networkName, region, vlanId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateVrack(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "name", networkName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "vlan_id", fmt.Sprintf("%d", vlanId)),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "id"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "current_state.vlan_id", fmt.Sprintf("%d", vlanId)),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_network_private_vrack.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudNetworkPrivateVrackImportStateIdFunc("ovh_cloud_network_private_vrack.test"),
			},
		},
	})
}

func TestAccCloudNetworkPrivateVrack_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	networkName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}
`, serviceName, networkName, region)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}
`, serviceName, updatedName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateVrack(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "name", networkName),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "checksum"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "name", updatedName),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "checksum"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "resource_status", "READY"),
				),
			},
		},
	})
}

// TestAccCloudNetworkPrivateVrack_serviceNameFromEnv validates that when service_name
// is omitted from the network and subnet configurations, the provider falls back to the
// OVH_CLOUD_PROJECT_SERVICE environment variable, and that the following plans neither
// show a diff nor replace the resources when another attribute changes.
func TestAccCloudNetworkPrivateVrack_serviceNameFromEnv(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	// Make the project id resolvable from the environment so the resources can be
	// configured without an explicit service_name.
	t.Setenv("OVH_CLOUD_PROJECT_SERVICE", serviceName)

	networkName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackNamePrefix)

	config := func(subnetName string) string {
		return fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "test" {
  name   = "%s"
  region = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "test" {
  network_id = ovh_cloud_network_private_vrack.test.id
  name       = "%s"
  cidr       = "10.0.0.0/24"
  region     = "%s"
}
`, networkName, region, subnetName, region)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateVrack(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(networkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "resource_status", "READY"),
				),
			},
			{
				// Re-planning with service_name still omitted must be a no-op.
				Config: config(networkName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				// Renaming the subnet updates it in place: service_name is kept from the
				// state instead of becoming unknown and triggering RequiresReplace.
				Config: config(networkName + "-renamed"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("ovh_cloud_network_private_vrack_subnet.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("ovh_cloud_network_private_vrack.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "service_name", serviceName),
			},
		},
	})
}

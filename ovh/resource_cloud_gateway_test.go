package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceCloudGatewayNamePrefix = "tf-test-gateway-v2-"

func TestAccCloudGateway_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	gatewayModel := os.Getenv("OVH_CLOUD_GATEWAY_MODEL_TEST")
	if gatewayModel == "" {
		gatewayModel = "S"
	}

	gatewayName := acctest.RandomWithPrefix(testAccResourceCloudGatewayNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_gateway" "test" {
  service_name = "%s"
  name         = "%s"

  region       = "%s"

  external_gateway = {
    enabled = true
    model   = "%s"
  }
}
`, serviceName, gatewayName, region, gatewayModel)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudGateway(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "name", gatewayName),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "external_gateway.enabled", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "external_gateway.model", gatewayModel),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_gateway.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_gateway.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_gateway.test", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_gateway.test", "current_state.name"),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "current_state.external_gateway.enabled", "true"),
					resource.TestCheckResourceAttrSet("ovh_cloud_gateway.test", "current_state.status"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_gateway.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudGatewayImportStateIdFunc("ovh_cloud_gateway.test"),
			},
		},
	})
}

func TestAccCloudGateway_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	gatewayModel := os.Getenv("OVH_CLOUD_GATEWAY_MODEL_TEST")
	if gatewayModel == "" {
		gatewayModel = "S"
	}

	gatewayName := acctest.RandomWithPrefix(testAccResourceCloudGatewayNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudGatewayNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_gateway" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  description  = "initial description"

  external_gateway = {
    enabled = true
    model   = "%s"
  }
}
`, serviceName, gatewayName, region, gatewayModel)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_gateway" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  description  = "updated description"

  external_gateway = {
    enabled = true
    model   = "%s"
  }
}
`, serviceName, updatedName, region, gatewayModel)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudGateway(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "name", gatewayName),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "description", "initial description"),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "external_gateway.enabled", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "external_gateway.model", gatewayModel),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "name", updatedName),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "description", "updated description"),
					resource.TestCheckResourceAttrSet("ovh_cloud_gateway.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_gateway.test", "checksum"),
				),
			},
		},
	})
}

func TestAccCloudGateway_disabledExternalGateway(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	gatewayName := acctest.RandomWithPrefix(testAccResourceCloudGatewayNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_gateway" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"

  external_gateway = {
    enabled = false
  }
}
`, serviceName, gatewayName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudGateway(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "name", gatewayName),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "external_gateway.enabled", "false"),
					resource.TestCheckResourceAttrSet("ovh_cloud_gateway.test", "id"),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "current_state.external_gateway.enabled", "false"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_gateway.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudGatewayImportStateIdFunc("ovh_cloud_gateway.test"),
			},
		},
	})
}

func TestAccCloudGateway_updateExternalGatewayModel(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	gatewayName := acctest.RandomWithPrefix(testAccResourceCloudGatewayNamePrefix)

	configTemplate := `
resource "ovh_cloud_gateway" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"

  external_gateway = {
    enabled = true
    model   = "%s"
  }
}
`

	config := fmt.Sprintf(configTemplate, serviceName, gatewayName, region, "S")
	updatedConfig := fmt.Sprintf(configTemplate, serviceName, gatewayName, region, "M")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudGateway(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "external_gateway.model", "S"),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "current_state.external_gateway.model", "S"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "external_gateway.model", "M"),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "current_state.external_gateway.model", "M"),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "resource_status", "READY"),
				),
			},
		},
	})
}

// TestAccCloudGateway_withSubnets exercises the full networking feature end to end:
// a private network, a subnet inside it, and a gateway attaching that subnet through
// subnet_ids (router interface).
func TestAccCloudGateway_withSubnets(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	gatewayModel := os.Getenv("OVH_CLOUD_GATEWAY_MODEL_TEST")
	if gatewayModel == "" {
		gatewayModel = "S"
	}

	networkName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackNamePrefix)
	subnetName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateSubnetNamePrefix)
	gatewayName := acctest.RandomWithPrefix(testAccResourceCloudGatewayNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "subnet" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id   = ovh_cloud_network_private_vrack.network.id
  name         = "%s"
  cidr         = "10.0.0.0/24"
  gateway_ip   = "10.0.0.1"
  dhcp_enabled = true
  region       = "%s"
}

resource "ovh_cloud_gateway" "test" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  name         = "%s"
  region       = "%s"
  subnet_ids   = [ovh_cloud_network_private_vrack_subnet.subnet.id]

  external_gateway = {
    enabled = true
    model   = "%s"
  }
}
`, serviceName, networkName, region, subnetName, region, gatewayName, region, gatewayModel)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudGateway(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "name", gatewayName),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "subnet_ids.#", "1"),
					resource.TestCheckResourceAttrPair(
						"ovh_cloud_gateway.test", "subnet_ids.0",
						"ovh_cloud_network_private_vrack_subnet.subnet", "id",
					),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "current_state.subnets.#", "1"),
					resource.TestCheckResourceAttrSet("ovh_cloud_gateway.test", "current_state.external_ip"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_gateway.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"subnet_ids",
				},
				ImportStateIdFunc: testAccCloudGatewayImportStateIdFunc("ovh_cloud_gateway.test"),
			},
		},
	})
}

func testAccPreCheckCloudGateway(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST not set")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_REGION_TEST not set")
	}
}

func testAccCloudGatewayImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}

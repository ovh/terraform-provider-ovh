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

  external_gateway {
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

  external_gateway {
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

  external_gateway {
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

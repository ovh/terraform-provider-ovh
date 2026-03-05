package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceCloudFloatingIpNamePrefix = "tf-test-fip-v2-"

func testAccPreCheckCloudFloatingIp(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST must be set for acceptance tests")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_REGION_TEST must be set for acceptance tests")
	}
}

func testAccCloudFloatingIpImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}

func TestAccCloudFloatingIp_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	description := testAccResourceCloudFloatingIpNamePrefix + "basic"

	config := fmt.Sprintf(`
resource "ovh_cloud_floating_ip" "test" {
  service_name = "%s"
  region       = "%s"
  description  = "%s"
}
`, serviceName, region, description)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudFloatingIp(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "description", description),
					resource.TestCheckResourceAttrSet("ovh_cloud_floating_ip.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_floating_ip.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_floating_ip.test", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_floating_ip.test", "current_state.ip"),
					resource.TestCheckResourceAttrSet("ovh_cloud_floating_ip.test", "current_state.status"),
					resource.TestCheckResourceAttrSet("ovh_cloud_floating_ip.test", "current_state.network_id"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_floating_ip.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudFloatingIpImportStateIdFunc("ovh_cloud_floating_ip.test"),
			},
		},
	})
}

func TestAccCloudFloatingIp_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	description := testAccResourceCloudFloatingIpNamePrefix + "before-update"
	updatedDescription := testAccResourceCloudFloatingIpNamePrefix + "after-update"

	config := fmt.Sprintf(`
resource "ovh_cloud_floating_ip" "test" {
  service_name = "%s"
  region       = "%s"
  description  = "%s"
}
`, serviceName, region, description)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_floating_ip" "test" {
  service_name = "%s"
  region       = "%s"
  description  = "%s"
}
`, serviceName, region, updatedDescription)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudFloatingIp(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "description", description),
					resource.TestCheckResourceAttrSet("ovh_cloud_floating_ip.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_floating_ip.test", "checksum"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "description", updatedDescription),
					resource.TestCheckResourceAttrSet("ovh_cloud_floating_ip.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_floating_ip.test", "checksum"),
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "resource_status", "READY"),
				),
			},
		},
	})
}

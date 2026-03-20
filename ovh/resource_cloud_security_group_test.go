package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceCloudSecurityGroupNamePrefix = "tf-test-sg-v2-"

func testAccPreCheckCloudSecurityGroup(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST must be set for acceptance tests")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_REGION_TEST must be set for acceptance tests")
	}
}

func testAccCloudSecurityGroupImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}

func TestAccCloudSecurityGroup_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	name := testAccResourceCloudSecurityGroupNamePrefix + "basic"

	config := fmt.Sprintf(`
resource "ovh_cloud_security_group" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  description  = "Test security group"

  rule {
    direction        = "INGRESS"
    ethernet_type       = "IPV4"
    protocol         = "TCP"
    port_range_min   = 22
    port_range_max   = 22
    remote_ip_prefix = "0.0.0.0/0"
    description      = "SSH"
  }
}
`, serviceName, region, name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudSecurityGroup(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "description", "Test security group"),
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "rule.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "rule.0.direction", "INGRESS"),
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "rule.0.ethernet_type", "IPV4"),
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "rule.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "rule.0.port_range_min", "22"),
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "rule.0.port_range_max", "22"),
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "rule.0.remote_ip_prefix", "0.0.0.0/0"),
					resource.TestCheckResourceAttrSet("ovh_cloud_security_group.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_security_group.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_security_group.test", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "resource_status", "READY"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_security_group.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudSecurityGroupImportStateIdFunc("ovh_cloud_security_group.test"),
			},
		},
	})
}

func TestAccCloudSecurityGroup_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	name := testAccResourceCloudSecurityGroupNamePrefix + "before-update"
	updatedName := testAccResourceCloudSecurityGroupNamePrefix + "after-update"

	config := fmt.Sprintf(`
resource "ovh_cloud_security_group" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  description  = "Before update"

  rule {
    direction        = "INGRESS"
    ethernet_type       = "IPV4"
    protocol         = "TCP"
    port_range_min   = 22
    port_range_max   = 22
    remote_ip_prefix = "0.0.0.0/0"
    description      = "SSH"
  }
}
`, serviceName, region, name)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_security_group" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  description  = "After update"

  rule {
    direction        = "INGRESS"
    ethernet_type       = "IPV4"
    protocol         = "TCP"
    port_range_min   = 22
    port_range_max   = 22
    remote_ip_prefix = "0.0.0.0/0"
    description      = "SSH"
  }

  rule {
    direction        = "INGRESS"
    ethernet_type       = "IPV4"
    protocol         = "TCP"
    port_range_min   = 443
    port_range_max   = 443
    remote_ip_prefix = "0.0.0.0/0"
    description      = "HTTPS"
  }
}
`, serviceName, region, updatedName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudSecurityGroup(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "description", "Before update"),
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "rule.#", "1"),
					resource.TestCheckResourceAttrSet("ovh_cloud_security_group.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_security_group.test", "checksum"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "name", updatedName),
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "description", "After update"),
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "rule.#", "2"),
					resource.TestCheckResourceAttrSet("ovh_cloud_security_group.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_security_group.test", "checksum"),
					resource.TestCheckResourceAttr("ovh_cloud_security_group.test", "resource_status", "READY"),
				),
			},
		},
	})
}

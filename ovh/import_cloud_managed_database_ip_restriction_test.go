package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudManagedDatabaseIpRestriction_importBasic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	engine := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_ENGINE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	ip := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_IP_RESTRICTION_IP_TEST")
	description := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccCloudManagedDatabaseIpRestrictionConfig,
		serviceName,
		description,
		engine,
		version,
		region,
		flavor,
		ip,
		description,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseIpRestriction(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "ovh_cloud_managed_database_ip_restriction.ip",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudManagedDatabaseIpRestrictionImportId("ovh_cloud_managed_database_ip_restriction.ip"),
			},
		},
	})
}

func testAccCloudManagedDatabaseIpRestrictionImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testIpRestriction, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_cloud_managed_database_ip_restriction not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s/%s",
			testIpRestriction.Primary.Attributes["service_name"],
			testIpRestriction.Primary.Attributes["engine"],
			testIpRestriction.Primary.Attributes["cluster_id"],
			testIpRestriction.Primary.Attributes["ip"],
		), nil
	}
}

package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudManagedDatabaseM3dbNamespace_importBasic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_M3DB_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	name := "mynamespace"
	resolution := "P2D"

	config := fmt.Sprintf(
		testAccCloudManagedDatabaseM3dbNamespaceConfig_basic,
		serviceName,
		description,
		version,
		region,
		flavor,
		name,
		resolution,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "ovh_cloud_managed_database_m3db_namespace.namespace",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudManagedDatabaseM3dbNamespaceImportId("ovh_cloud_managed_database_m3db_namespace.namespace"),
			},
		},
	})
}

func testAccCloudManagedDatabaseM3dbNamespaceImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testM3dbNamespace, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_cloud_managed_database_m3db_namespace not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s",
			testM3dbNamespace.Primary.Attributes["service_name"],
			testM3dbNamespace.Primary.Attributes["cluster_id"],
			testM3dbNamespace.Primary.Attributes["id"],
		), nil
	}
}

package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudProjectDatabaseM3dbUser_importBasic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_M3DB_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_M3DB_FLAVOR_TEST")
	if flavor == "" {
		flavor = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	}
	description := acctest.RandomWithPrefix(test_prefix)
	name := "johndoe"
	group := "mygroup"

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseM3dbUserConfig_basic,
		serviceName,
		description,
		version,
		region,
		flavor,
		name,
		group,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            "ovh_cloud_project_database_m3db_user.user",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccCloudProjectDatabaseM3dbUserImportId("ovh_cloud_project_database_m3db_user.user"),
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccCloudProjectDatabaseM3dbUserImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testM3dbUser, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_cloud_project_database_m3db_user not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s",
			testM3dbUser.Primary.Attributes["service_name"],
			testM3dbUser.Primary.Attributes["cluster_id"],
			testM3dbUser.Primary.Attributes["id"],
		), nil
	}
}

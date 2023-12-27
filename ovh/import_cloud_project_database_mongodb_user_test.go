package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCloudProjectDatabaseMongodbUser_importBasic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_MONGODB_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_MONGODB_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_MONGODB_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	name := "johndoe"
	rolesBackup := "backup"
	rolesReadAnyDatabase := "readAnyDatabase"

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseMongodbUserConfig_basic,
		serviceName,
		description,
		version,
		region,
		region,
		region,
		flavor,
		name,
		rolesBackup,
		rolesReadAnyDatabase,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseMongoDBNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            "ovh_cloud_project_database_mongodb_user.user",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccCloudProjectDatabaseMongodbUserImportId("ovh_cloud_project_database_mongodb_user.user"),
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccCloudProjectDatabaseMongodbUserImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testMongodbUser, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_cloud_project_database_mongodb_user not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s",
			testMongodbUser.Primary.Attributes["service_name"],
			testMongodbUser.Primary.Attributes["cluster_id"],
			testMongodbUser.Primary.Attributes["id"],
		), nil
	}
}

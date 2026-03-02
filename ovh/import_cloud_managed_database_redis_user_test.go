package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudManagedDatabaseRedisUser_importBasic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REDIS_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	categoriesSet := "+@set"
	categoriesSortedset := "+@sortedset"
	channels := "*"
	commandsGet := "+get"
	commandsSet := "-set"
	keysData := "data"
	keysProperties := "properties"
	name := "johndoe"

	config := fmt.Sprintf(
		testAccCloudManagedDatabaseRedisUserConfig,
		serviceName,
		description,
		version,
		region,
		flavor,
		categoriesSet,
		categoriesSortedset,
		channels,
		commandsGet,
		commandsSet,
		keysData,
		keysProperties,
		name,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            "ovh_cloud_managed_database_redis_user.user",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccCloudManagedDatabaseRedisUserImportId("ovh_cloud_managed_database_redis_user.user"),
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccCloudManagedDatabaseRedisUserImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testRedisUser, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_cloud_managed_database_redis_user not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s",
			testRedisUser.Primary.Attributes["service_name"],
			testRedisUser.Primary.Attributes["cluster_id"],
			testRedisUser.Primary.Attributes["id"],
		), nil
	}
}

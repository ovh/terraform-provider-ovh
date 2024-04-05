package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudProjectDatabaseLogSubscription_importBasic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	engine := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_ENGINE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	LDPserviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	title := acctest.RandomWithPrefix(test_prefix)
	desc := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseLogSubscrition_basic,
		serviceName,
		description,
		engine,
		version,
		region,
		flavor,
		LDPserviceName,
		title,
		desc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckDbaasLogs(t)
			testAccPreCheckCloudDatabase(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            "ovh_cloud_project_database_log_subscription.sub",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccCloudProjectDatabaseLogSubscriptionImportID("ovh_cloud_project_database_log_subscription.sub"),
				ImportStateVerifyIgnore: []string{"operation_id"},
			},
		},
	})
}

func testAccCloudProjectDatabaseLogSubscriptionImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testUser, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_cloud_project_database_log_subscription not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s/%s",
			testUser.Primary.Attributes["service_name"],
			testUser.Primary.Attributes["engine"],
			testUser.Primary.Attributes["cluster_id"],
			testUser.Primary.Attributes["id"],
		), nil
	}
}

package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudProjectDatabaseKafkaAcl_importBasic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_KAFKA_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	permission := "read"
	topic := "myTopic"
	username := "johnDoe"

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseKafkaAclConfig,
		serviceName,
		description,
		version,
		region,
		region,
		region,
		flavor,
		permission,
		topic,
		username,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "ovh_cloud_project_database_kafka_acl.acl",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudProjectDatabaseKafkaAclImportId("ovh_cloud_project_database_kafka_acl.acl"),
			},
		},
	})
}

func testAccCloudProjectDatabaseKafkaAclImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testKafkaAcl, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_cloud_project_database_kafka_acl not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s",
			testKafkaAcl.Primary.Attributes["service_name"],
			testKafkaAcl.Primary.Attributes["cluster_id"],
			testKafkaAcl.Primary.Attributes["id"],
		), nil
	}
}

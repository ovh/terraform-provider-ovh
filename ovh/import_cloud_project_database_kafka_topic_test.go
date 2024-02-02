package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudProjectDatabaseKafkaTopic_importBasic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_KAFKA_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	name1 := "myTopic1"
	name2 := "myTopic2"

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseKafkaTopicConfig,
		serviceName,
		description,
		version,
		region,
		region,
		region,
		flavor,
		name1,
		name2,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "ovh_cloud_project_database_kafka_topic.topic1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudProjectDatabaseKafkaTopicImportId("ovh_cloud_project_database_kafka_topic.topic1"),
			},
		},
	})
}

func testAccCloudProjectDatabaseKafkaTopicImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testKafkaTopic, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_cloud_project_database_kafka_topic not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s",
			testKafkaTopic.Primary.Attributes["service_name"],
			testKafkaTopic.Primary.Attributes["cluster_id"],
			testKafkaTopic.Primary.Attributes["id"],
		), nil
	}
}

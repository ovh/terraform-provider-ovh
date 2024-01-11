package ovh

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCloudProjectDatabaseKafkaTopicDatasourceConfig_Basic = `
resource "ovh_cloud_project_database" "db" {
	service_name = "%s"
	description  = "%s"
	engine       = "kafka"
	version      = "%s"
	plan         = "business"
	nodes {
		region     = "%s"
	}
	nodes {
		region     = "%s"
	}
	nodes {
		region     = "%s"
	}
	flavor = "%s"
}

resource "ovh_cloud_project_database_kafka_topic" "topic" {
	service_name = ovh_cloud_project_database.db.service_name
	cluster_id   = ovh_cloud_project_database.db.id
	name = "%s"
}

data "ovh_cloud_project_database_kafka_topic" "topic" {
  service_name = ovh_cloud_project_database_kafka_topic.topic.service_name
  cluster_id   = ovh_cloud_project_database_kafka_topic.topic.cluster_id
  id           = ovh_cloud_project_database_kafka_topic.topic.id
}
`

func TestAccCloudProjectDatabaseKafkaTopicDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_KAFKA_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	name := "myTopic"

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseKafkaTopicDatasourceConfig_Basic,
		serviceName,
		description,
		version,
		region,
		region,
		region,
		flavor,
		name,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_kafka_topic.topic", "name", name),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_kafka_topic.topic", "min_insync_replicas", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_kafka_topic.topic", "partitions", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_kafka_topic.topic", "replication", strconv.Itoa(2)),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_kafka_topic.topic", "retention_bytes", strconv.Itoa(-1)),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_kafka_topic.topic", "retention_hours", strconv.Itoa(168)),
				),
			},
		},
	})
}

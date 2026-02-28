package ovh

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCloudManagedAnalyticsKafkaTopicConfig = `
resource "ovh_cloud_managed_analytics" "db" {
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

resource "ovh_cloud_managed_analytics_kafka_topic" "topic1" {
	service_name = ovh_cloud_project_database.db.service_name
	cluster_id   = ovh_cloud_project_database.db.id
	name = "%s"
}

resource "ovh_cloud_managed_analytics_kafka_topic" "topic2" {
	service_name = ovh_cloud_project_database.db.service_name
	cluster_id   = ovh_cloud_project_database.db.id
	name = "%s"
	min_insync_replicas = 1
	partitions = 3
	replication = 2
	retention_bytes = 4
	retention_hours = 5
}
`

func TestAccCloudManagedAnalyticsKafkaTopic_basic(t *testing.T) {
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
		testAccCloudManagedAnalyticsKafkaTopicConfig,
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_analytics_kafka_topic.topic1", "name", name1),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_analytics_kafka_topic.topic1", "min_insync_replicas", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_analytics_kafka_topic.topic1", "partitions", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_analytics_kafka_topic.topic1", "replication", strconv.Itoa(2)),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_analytics_kafka_topic.topic1", "retention_bytes", strconv.Itoa(-1)),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_analytics_kafka_topic.topic1", "retention_hours", strconv.Itoa(168)),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_analytics_kafka_topic.topic2", "name", name2),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_analytics_kafka_topic.topic2", "min_insync_replicas", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_analytics_kafka_topic.topic2", "partitions", strconv.Itoa(3)),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_analytics_kafka_topic.topic2", "replication", strconv.Itoa(2)),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_analytics_kafka_topic.topic2", "retention_bytes", strconv.Itoa(4)),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_analytics_kafka_topic.topic2", "retention_hours", strconv.Itoa(5)),
				),
			},
		},
	})
}

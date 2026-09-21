package ovh

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

const testAccCloudProjectDatabaseKafkaTopicConfig = `
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

resource "ovh_cloud_project_database_kafka_topic" "topic1" {
	service_name = ovh_cloud_project_database.db.service_name
	cluster_id   = ovh_cloud_project_database.db.id
	name = "%s"
}

resource "ovh_cloud_project_database_kafka_topic" "topic2" {
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

func TestAccCloudProjectDatabaseKafkaTopic_basic(t *testing.T) {
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_kafka_topic.topic1", "name", name1),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_kafka_topic.topic1", "min_insync_replicas", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_kafka_topic.topic1", "partitions", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_kafka_topic.topic1", "replication", strconv.Itoa(2)),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_kafka_topic.topic1", "retention_bytes", strconv.Itoa(-1)),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_kafka_topic.topic1", "retention_hours", strconv.Itoa(168)),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_kafka_topic.topic2", "name", name2),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_kafka_topic.topic2", "min_insync_replicas", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_kafka_topic.topic2", "partitions", strconv.Itoa(3)),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_kafka_topic.topic2", "replication", strconv.Itoa(2)),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_kafka_topic.topic2", "retention_bytes", strconv.Itoa(4)),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_kafka_topic.topic2", "retention_hours", strconv.Itoa(5)),
				),
			},
		},
	})
}

const testAccCloudProjectDatabaseKafkaTopicUpdateConfig = `
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
	min_insync_replicas = %d
	partitions = %d
	replication = %d
	retention_bytes = %d
	retention_hours = %d
}
`

// captureKafkaTopicID stores the topic id on first call, then fails on any
// later call where it changed, proving the topic was updated and not recreated.
func captureKafkaTopicID(rn string, store *string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(rn, "id", func(v string) error {
		if v == "" {
			return fmt.Errorf("expected topic id to be set")
		}
		if *store == "" {
			*store = v
			return nil
		}
		if v != *store {
			return fmt.Errorf("topic was replaced: id changed from %q to %q", *store, v)
		}
		return nil
	})
}

func TestAccCloudProjectDatabaseKafkaTopic_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_KAFKA_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	name := "myTopic"
	rn := "ovh_cloud_project_database_kafka_topic.topic"

	config := func(minInsyncReplicas, partitions, replication, retentionBytes, retentionHours int) string {
		return fmt.Sprintf(
			testAccCloudProjectDatabaseKafkaTopicUpdateConfig,
			serviceName,
			description,
			version,
			region,
			region,
			region,
			flavor,
			name,
			minInsyncReplicas,
			partitions,
			replication,
			retentionBytes,
			retentionHours,
		)
	}

	checkAttrs := func(minInsyncReplicas, partitions, replication, retentionBytes, retentionHours int) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(rn, "name", name),
			resource.TestCheckResourceAttr(rn, "min_insync_replicas", strconv.Itoa(minInsyncReplicas)),
			resource.TestCheckResourceAttr(rn, "partitions", strconv.Itoa(partitions)),
			resource.TestCheckResourceAttr(rn, "replication", strconv.Itoa(replication)),
			resource.TestCheckResourceAttr(rn, "retention_bytes", strconv.Itoa(retentionBytes)),
			resource.TestCheckResourceAttr(rn, "retention_hours", strconv.Itoa(retentionHours)),
		)
	}

	var topicID string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config(1, 3, 2, 4, 5),
				Check: resource.ComposeTestCheckFunc(
					checkAttrs(1, 3, 2, 4, 5),
					captureKafkaTopicID(rn, &topicID),
				),
			},
			{
				// min_insync_replicas, replication and both retentions update in place.
				Config: config(2, 3, 3, 8, 10),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					checkAttrs(2, 3, 3, 8, 10),
					captureKafkaTopicID(rn, &topicID),
				),
			},
			{
				// Partitions can be increased in place.
				Config: config(2, 6, 3, 8, 10),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					checkAttrs(2, 6, 3, 8, 10),
					captureKafkaTopicID(rn, &topicID),
				),
			},
			{
				// The API rejects a partition decrease, so it must force a new resource.
				Config: config(2, 3, 3, 8, 10),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionReplace),
					},
				},
				Check: checkAttrs(2, 3, 3, 8, 10),
			},
		},
	})
}

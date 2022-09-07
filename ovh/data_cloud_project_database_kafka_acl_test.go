package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccCloudProjectDatabaseKafkaAclDatasourceConfig_Basic = `
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

resource "ovh_cloud_project_database_kafka_acl" "acl" {
	service_name = ovh_cloud_project_database.db.service_name
	cluster_id   = ovh_cloud_project_database.db.id
	permission	 = "%s"
	topic 		 = "%s"
	username 	 = "%s"
}

data "ovh_cloud_project_database_kafka_acl" "acl" {
  service_name = ovh_cloud_project_database_kafka_acl.acl.service_name
  cluster_id   = ovh_cloud_project_database_kafka_acl.acl.cluster_id
  id           = ovh_cloud_project_database_kafka_acl.acl.id
}
`

func TestAccCloudProjectDatabaseKafkaAclDataSource_basic(t *testing.T) {
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
		testAccCloudProjectDatabaseKafkaAclDatasourceConfig_Basic,
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_kafka_acl.acl", "permission", permission),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_kafka_acl.acl", "topic", topic),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_kafka_acl.acl", "username", username),
				),
			},
		},
	})
}

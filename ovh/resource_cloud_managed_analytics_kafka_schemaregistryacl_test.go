package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCloudManagedAnalyticsKafkaSchemaRegistryAclConfig = `
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

resource "ovh_cloud_managed_analytics_kafka_schemaregistryacl" "schemaRegistryAcl" {
	service_name = ovh_cloud_project_database.db.service_name
	cluster_id   = ovh_cloud_project_database.db.id
	permission	 = "%s"
	resource 	 = "%s"
	username 	 = "%s"
}
`

func TestAccCloudManagedAnalyticsKafkaSchemaRegistryAcl_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_KAFKA_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	permission := "schema_registry_read"
	aclResource := "Subject:myResource"
	username := "johnDoe"

	config := fmt.Sprintf(
		testAccCloudManagedAnalyticsKafkaSchemaRegistryAclConfig,
		serviceName,
		description,
		version,
		region,
		region,
		region,
		flavor,
		permission,
		aclResource,
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
						"ovh_cloud_managed_analytics_kafka_schemaregistryacl.schemaRegistryAcl", "permission", permission),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_analytics_kafka_schemaregistryacl.schemaRegistryAcl", "resource", aclResource),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_analytics_kafka_schemaregistryacl.schemaRegistryAcl", "username", username),
				),
			},
		},
	})
}

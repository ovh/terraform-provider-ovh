package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudProjectDatabaseKafkaSchemaRegistryAcl_importBasic(t *testing.T) {
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
		testAccCloudProjectDatabaseKafkaSchemaRegistryAclConfig,
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
			},
			{
				ResourceName:      "ovh_cloud_project_database_kafka_schemaregistryacl.schemaRegistryAcl",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudProjectDatabaseKafkaSchemaRegistryAclImportId("ovh_cloud_project_database_kafka_schemaregistryacl.schemaRegistryAcl"),
			},
		},
	})
}

func testAccCloudProjectDatabaseKafkaSchemaRegistryAclImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testKafkaSchemaRegistryAcl, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_cloud_project_database_kafka_schemaregistryacl not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s",
			testKafkaSchemaRegistryAcl.Primary.Attributes["service_name"],
			testKafkaSchemaRegistryAcl.Primary.Attributes["cluster_id"],
			testKafkaSchemaRegistryAcl.Primary.Attributes["id"],
		), nil
	}
}

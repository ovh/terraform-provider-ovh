package ovh

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccCloudProjectDatabaseOpensearchPatternDatasourceConfig_Basic = `
resource "ovh_cloud_project_database" "db" {
	service_name = "%s"
	description  = "%s"
	engine       = "opensearch"
	version      = "%s"
	plan         = "essential"
	nodes {
		region     = "%s"
	}
	flavor = "%s"
}

resource "ovh_cloud_project_database_opensearch_pattern" "pattern" {
	service_name = ovh_cloud_project_database.db.service_name
	cluster_id   = ovh_cloud_project_database.db.id
	max_index_count = %d
	pattern = "%s"
}

data "ovh_cloud_project_database_opensearch_pattern" "pattern" {
  service_name = ovh_cloud_project_database_opensearch_pattern.pattern.service_name
  cluster_id   = ovh_cloud_project_database_opensearch_pattern.pattern.cluster_id
  id           = ovh_cloud_project_database_opensearch_pattern.pattern.id
}
`

func TestAccCloudProjectDatabaseOpensearchPatternDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_OPENSEARCH_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	maxIndexCount := 2
	pattern := "logs_*"

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseOpensearchPatternDatasourceConfig_Basic,
		serviceName,
		description,
		version,
		region,
		flavor,
		maxIndexCount,
		pattern,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_opensearch_pattern.pattern", "max_index_count", strconv.Itoa(maxIndexCount)),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_opensearch_pattern.pattern", "pattern", pattern),
				),
			},
		},
	})
}

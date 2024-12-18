package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCloudProjectDatabaseMongodbPrometheusConfig = `
resource "ovh_cloud_project_database" "db" {
	service_name = "%s"
	description  = "%s"
	engine       = "mongodb"
	version      = "%s"
	plan         = "production"
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

resource "ovh_cloud_project_database_mongodb_prometheus" "prometheus" {
	service_name = ovh_cloud_project_database.db.service_name
	cluster_id   = ovh_cloud_project_database.db.id
}
`

func TestAccCloudProjectDatabaseMongodbPrometheus_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_MONGODB_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_MONGODB_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_MONGODB_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseMongodbPrometheusConfig,
		serviceName,
		description,
		version,
		region,
		region,
		region,
		flavor,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabase(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database_mongodb_prometheus.prometheus", "password"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_mongodb_prometheus.prometheus", "username", "prometheus"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database_mongodb_prometheus.prometheus", "srv_domain"),
				),
			},
		},
	})
}

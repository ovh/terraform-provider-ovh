package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCloudManagedDatabasePrometheusConfig = `
resource "ovh_cloud_managed_database" "db" {
	service_name = "%s"
	description  = "%s"
	engine       = "%s"
	version      = "%s"
	plan         = "essential"
	nodes {
		region     = "%s"
	}
	flavor = "%s"
}

resource "ovh_cloud_managed_database_prometheus" "prometheus" {
	service_name = ovh_cloud_project_database.db.service_name
	engine		 = ovh_cloud_project_database.db.engine
	cluster_id   = ovh_cloud_project_database.db.id
}
`

func TestAccCloudManagedDatabasePrometheus_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	engine := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_ENGINE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccCloudManagedDatabasePrometheusConfig,
		serviceName,
		description,
		engine,
		version,
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
						"ovh_cloud_managed_database_prometheus.prometheus", "password"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database_prometheus.prometheus", "targets.#"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database_prometheus.prometheus", "targets.0.host"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database_prometheus.prometheus", "targets.0.port"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_database_prometheus.prometheus", "username", "prometheus"),
				),
			},
		},
	})
}

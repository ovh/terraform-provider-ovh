package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCloudManagedDatabaseIpRestrictionsDatasourceConfig_Basic = `
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

resource "ovh_cloud_managed_database_ip_restriction" "ip" {
	service_name = ovh_cloud_project_database.db.service_name
	engine		 = ovh_cloud_project_database.db.engine
	cluster_id   = ovh_cloud_project_database.db.id
	ip			 = "%s"
}

data "ovh_cloud_managed_database_ip_restrictions" "ips" {
  service_name = ovh_cloud_project_database_ip_restriction.ip.service_name
  engine	   = ovh_cloud_project_database_ip_restriction.ip.engine
  cluster_id   = ovh_cloud_project_database_ip_restriction.ip.cluster_id
}
`

func TestAccCloudManagedDatabaseIpRestrictionsDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	engine := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_ENGINE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	ip := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_IP_RESTRICTION_IP_TEST")
	description := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccCloudManagedDatabaseIpRestrictionsDatasourceConfig_Basic,
		serviceName,
		description,
		engine,
		version,
		region,
		flavor,
		ip,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudDatabaseIpRestriction(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_ip_restrictions.ips",
						"ips.#",
					),
				),
			},
		},
	})
}

package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCloudManagedDatabaseM3dbUserDatasourceConfig_Basic = `
resource "ovh_cloud_managed_database" "db" {
	service_name = "%s"
	description  = "%s"
	engine       = "m3db"
	version      = "%s"
	plan         = "essential"
	nodes {
		region     = "%s"
	}
	flavor = "%s"
}

resource "ovh_cloud_managed_database_m3db_user" "user" {
	service_name = ovh_cloud_managed_database.db.service_name
	cluster_id   = ovh_cloud_managed_database.db.id
	name		 = "%s"
	group 		 = "%s"
}

data "ovh_cloud_managed_database_m3db_user" "user" {
  service_name = ovh_cloud_managed_database_m3db_user.user.service_name
  cluster_id   = ovh_cloud_managed_database_m3db_user.user.cluster_id
  name     = ovh_cloud_managed_database_m3db_user.user.name
}
`

func TestAccCloudManagedDatabaseM3dbUserDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_M3DB_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	name := "johndoe"
	group := "mygroup"

	config := fmt.Sprintf(
		testAccCloudManagedDatabaseM3dbUserDatasourceConfig_Basic,
		serviceName,
		description,
		version,
		region,
		flavor,
		name,
		group,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_managed_database_m3db_user.user", "created_at"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_managed_database_m3db_user.user", "name", name),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_database_m3db_user.user", "group", group),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_managed_database_m3db_user.user", "status"),
				),
			},
		},
	})
}

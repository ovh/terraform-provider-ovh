package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCloudProjectDatabaseUserDatasourceConfig_Basic = `
resource "ovh_cloud_project_database" "db" {
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

resource "ovh_cloud_project_database_user" "user" {
	service_name = ovh_cloud_project_database.db.service_name
	engine		 = ovh_cloud_project_database.db.engine
	cluster_id   = ovh_cloud_project_database.db.id
	name		 = "%s"
}

data "ovh_cloud_project_database_user" "user" {
  service_name = ovh_cloud_project_database_user.user.service_name
  engine	   = ovh_cloud_project_database_user.user.engine
  cluster_id   = ovh_cloud_project_database_user.user.cluster_id
  name 	   = ovh_cloud_project_database_user.user.name
}
`

func TestAccCloudProjectDatabaseUserDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	engine := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_ENGINE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	name := "johndoe"
	if engine == "grafana" {
		name = "avnadmin"
	}

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseUserDatasourceConfig_Basic,
		serviceName,
		description,
		engine,
		version,
		region,
		flavor,
		name,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabase(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_user.user", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_user.user", "status"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_user.user", "name", name,
					),
				),
			},
		},
	})
}

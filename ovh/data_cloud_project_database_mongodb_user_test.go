package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCloudProjectDatabaseMongodbUserDatasourceConfig_Basic = `
resource "ovh_cloud_project_database" "db" {
	service_name = "%s"
	description  = "%s"
	engine       = "mongodb"
	version      = "%s"
	plan         = "discovery"
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

resource "ovh_cloud_project_database_mongodb_user" "user" {
	service_name = ovh_cloud_project_database.db.service_name
	cluster_id   = ovh_cloud_project_database.db.id
	name		 = "%s"
	roles		 = ["%s", "%s"]
}

data "ovh_cloud_project_database_mongodb_user" "user" {
  service_name = ovh_cloud_project_database_mongodb_user.user.service_name
  cluster_id   = ovh_cloud_project_database_mongodb_user.user.cluster_id
  name         = ovh_cloud_project_database_mongodb_user.user.name
}
`

func TestAccCloudProjectDatabaseMongodbUserDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_MONGODB_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_MONGODB_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_MONGODB_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	name := "johndoe"
	rolesBackup := "backup@admin"
	rolesReadAnyDatabase := "readAnyDatabase@admin"

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseMongodbUserDatasourceConfig_Basic,
		serviceName,
		description,
		version,
		region,
		region,
		region,
		flavor,
		name,
		rolesBackup,
		rolesReadAnyDatabase,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseMongoDBNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_mongodb_user.user", "created_at"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_mongodb_user.user", "name", name+"@admin"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_mongodb_user.user", "roles.#"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_mongodb_user.user", "roles.0", rolesBackup),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_mongodb_user.user", "roles.1", rolesReadAnyDatabase),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_mongodb_user.user", "status"),
				),
			},
		},
	})
}

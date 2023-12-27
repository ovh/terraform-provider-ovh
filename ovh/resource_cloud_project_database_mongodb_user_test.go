package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccCloudProjectDatabaseMongodbUserConfig_basic = `
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
	roles = ["%s", "%s"]
}
`

func TestAccCloudProjectDatabaseMongodbUser_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_MONGODB_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_MONGODB_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_MONGODB_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	name := "johndoe"
	rolesBackup := "backup"
	rolesReadAnyDatabase := "readAnyDatabase"

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseMongodbUserConfig_basic,
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
						"ovh_cloud_project_database_mongodb_user.user", "created_at"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_mongodb_user.user", "name", name+"@admin"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database_mongodb_user.user", "password"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database_mongodb_user.user", "roles.#"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_mongodb_user.user", "roles.0", rolesBackup),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_mongodb_user.user", "roles.1", rolesReadAnyDatabase),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database_mongodb_user.user", "status"),
				),
			},
		},
	})
}

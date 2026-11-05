package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCloudProjectDatabasePostgresqlUserConfig = `
resource "ovh_cloud_project_database" "db" {
	service_name = "%s"
	description  = "%s"
	engine       = "postgresql"
	version      = "%s"
	plan         = "essential"
	nodes {
		region     = "%s"
	}
	flavor = "%s"
}

resource "ovh_cloud_project_database_postgresql_user" "user" {
	service_name = ovh_cloud_project_database.db.service_name
	cluster_id   = ovh_cloud_project_database.db.id
	name		 = "%s"
	roles 		 = ["%s"]
}
`

func TestAccCloudProjectDatabasePostgresqlUser_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_POSTGRESQL_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	name := "johndoe"
	replication := "replication"

	config := fmt.Sprintf(
		testAccCloudProjectDatabasePostgresqlUserConfig,
		serviceName,
		description,
		version,
		region,
		flavor,
		name,
		replication,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database_postgresql_user.user", "created_at"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database_postgresql_user.user", "password"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database_postgresql_user.user", "roles.#"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_postgresql_user.user", "roles.0", replication),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database_postgresql_user.user", "status"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_postgresql_user.user", "name", name,
					),
				),
			},
		},
	})
}

const testAccCloudProjectDatabasePostgresqlUserConfig_builtin = `
resource "ovh_cloud_project_database" "db" {
  service_name = "%s"
  description  = "%s"
  engine       = "postgresql"
  version      = "%s"
  plan         = "essential"
  nodes {
    region = "%s"
  }
  flavor = "%s"
}

data "ovh_cloud_project_database_postgresql_user" "builtin" {
  service_name = ovh_cloud_project_database.db.service_name
  cluster_id   = ovh_cloud_project_database.db.id
  name         = "avnadmin"
}

resource "ovh_cloud_project_database_postgresql_user" "builtin" {
  service_name   = ovh_cloud_project_database.db.service_name
  cluster_id     = ovh_cloud_project_database.db.id
  name           = "avnadmin"
  password_reset = "initial"
}
`

func TestAccCloudProjectDatabasePostgresqlUser_builtin(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_POSTGRESQL_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	config := fmt.Sprintf(
		testAccCloudProjectDatabasePostgresqlUserConfig_builtin,
		serviceName,
		acctest.RandomWithPrefix(test_prefix),
		version,
		region,
		flavor,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"ovh_cloud_project_database_postgresql_user.builtin", "id",
						"data.ovh_cloud_project_database_postgresql_user.builtin", "id"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database_postgresql_user.builtin", "id"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_postgresql_user.builtin", "name", "avnadmin"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database_postgresql_user.builtin", "password"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database_postgresql_user.builtin", "created_at"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_postgresql_user.builtin", "status", "READY"),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

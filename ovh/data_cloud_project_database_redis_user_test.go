package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccCloudProjectDatabaseRedisUserDatasourceConfig = `
resource "ovh_cloud_project_database" "db" {
	service_name = "%s"
	description  = "%s"
	engine       = "redis"
	version      = "%s"
	plan         = "essential"
	nodes {
		region     = "%s"
	}
	flavor = "%s"
}

resource "ovh_cloud_project_database_redis_user" "user" {
	service_name = ovh_cloud_project_database.db.service_name
	cluster_id   = ovh_cloud_project_database.db.id
	categories	 = ["%s", "%s"]
	channels	 = ["%s"]
	commands	 = ["%s", "%s"]
	keys		 = ["%s", "%s"]
	name		 = "%s"
}

data "ovh_cloud_project_database_redis_user" "user" {
  service_name = ovh_cloud_project_database_redis_user.user.service_name
  cluster_id   = ovh_cloud_project_database_redis_user.user.cluster_id
  name     = ovh_cloud_project_database_redis_user.user.name
}
`

func TestAccCloudProjectDatabaseRedisUserDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REDIS_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	categoriesSet := "+@set"
	categoriesSortedset := "+@sortedset"
	channels := "*"
	commandsGet := "+get"
	commandsSet := "-set"
	keysData := "data"
	keysProperties := "properties"
	name := "johndoe"

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseRedisUserDatasourceConfig,
		serviceName,
		description,
		version,
		region,
		flavor,
		categoriesSet,
		categoriesSortedset,
		channels,
		commandsGet,
		commandsSet,
		keysData,
		keysProperties,
		name,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_redis_user.user", "categories.#"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_redis_user.user", "categories.0", categoriesSet),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_redis_user.user", "categories.1", categoriesSortedset),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_redis_user.user", "channels.#"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_redis_user.user", "channels.0", channels),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_redis_user.user", "commands.#"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_redis_user.user", "commands.0", commandsGet),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_redis_user.user", "commands.1", commandsSet),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_redis_user.user", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_redis_user.user", "keys.#"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_redis_user.user", "keys.0", keysData),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_redis_user.user", "keys.1", keysProperties),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_redis_user.user", "name", name),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_redis_user.user", "status"),
				),
			},
		},
	})
}

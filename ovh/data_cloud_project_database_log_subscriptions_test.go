package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCloudProjectDatabaseLogSubsriptionsDatasourceConfig_Basic = `
resource "ovh_cloud_project_database" "db" {
	service_name = "%s"
	description  = "%s"
	engine       = "%s"
	version      = "%s"
	plan         = "essential"
	nodes {
		region = "%s"
	}
	flavor = "%s"
}

resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
	service_name = "%s"
	title        = "%s"
	description  = "%s"
}

resource "ovh_cloud_project_database_log_subscription" "sub" {
	service_name = ovh_cloud_project_database.db.service_name
	engine       = ovh_cloud_project_database.db.engine
	cluster_id   = ovh_cloud_project_database.db.id
	stream_id    = ovh_dbaas_logs_output_graylog_stream.stream.id
}

data "ovh_cloud_project_database_log_subscriptions" "subs" {
	service_name = ovh_cloud_project_database_log_subscription.sub.service_name
	engine       = ovh_cloud_project_database_log_subscription.sub.engine
	cluster_id   = ovh_cloud_project_database_log_subscription.sub.cluster_id
}
`

func TestAccCloudProjectDatabaseLogSubscriptionsDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	engine := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_ENGINE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	LDPserviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	title := acctest.RandomWithPrefix(test_prefix)
	desc := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseLogSubsriptionsDatasourceConfig_Basic,
		serviceName,
		description,
		engine,
		version,
		region,
		flavor,
		LDPserviceName,
		title,
		desc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckDbaasLogs(t)
			testAccPreCheckCloudDatabaseNoEngine(t)
		},

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_log_subscription.sub",
						"subscription_ids.#",
					),
				),
			},
		},
	})
}

package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccCloudProjectDatabasesDatasourceConfig_Basic = `
data "ovh_cloud_project_databases" "dbs" {
  service_name = "%s"
  engine = "%s"
}
`

func TestAccCloudProjectDatabasesDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	engine := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_ENGINE_TEST")
	config := fmt.Sprintf(
		testAccCloudProjectDatabasesDatasourceConfig_Basic,
		serviceName,
		engine,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabase(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_databases.dbs",
						"cluster_ids.#",
					),
				),
			},
		},
	})
}

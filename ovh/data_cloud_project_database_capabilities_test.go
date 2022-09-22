package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccCloudProjectDatabaseCapabilitiesDatasourceConfig_Basic = `
data "ovh_cloud_project_database_capabilities" "capabilities" {
	service_name = "%s"
}
`

func TestAccCloudProjectDatabaseCapabilitiesDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseCapabilitiesDatasourceConfig_Basic,
		serviceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_capabilities.capabilities",
						"engines.#",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_capabilities.capabilities",
						"engines.0.default_version",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_capabilities.capabilities",
						"engines.0.description",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_capabilities.capabilities",
						"engines.0.name",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_capabilities.capabilities",
						"engines.0.ssl_modes.#",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_capabilities.capabilities",
						"engines.0.versions.#",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_capabilities.capabilities",
						"flavors.#",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_capabilities.capabilities",
						"flavors.0.core",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_capabilities.capabilities",
						"flavors.0.memory",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_capabilities.capabilities",
						"flavors.0.name",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_capabilities.capabilities",
						"flavors.0.storage",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_capabilities.capabilities",
						"options.#",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_capabilities.capabilities",
						"options.0.name",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_capabilities.capabilities",
						"options.0.type",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_capabilities.capabilities",
						"plans.#",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_capabilities.capabilities",
						"plans.0.backup_retention",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_capabilities.capabilities",
						"plans.0.description",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_capabilities.capabilities",
						"plans.0.name",
					),
				),
			},
		},
	})
}

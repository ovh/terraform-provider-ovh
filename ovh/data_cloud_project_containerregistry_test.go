package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudProjectContainerRegistryDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	registryName := acctest.RandomWithPrefix(test_prefix)
	region := os.Getenv("OVH_CLOUD_PROJECT_CONTAINERREGISTRY_REGION_TEST")
	config := fmt.Sprintf(
		testAccCloudProjectContainerRegistryDatasourceConfig_Basic,
		serviceName,
		region,
		registryName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckContainerRegistry(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAccCloudProjectContainerRegistryDatasourceConfig_Basic = `
data "ovh_cloud_project_capabilities_containerregistry_filter" "regcap" {
	service_name = "%s"
    plan_name    = "SMALL"
    region       = "%s"
}

resource "ovh_cloud_project_containerregistry" "reg" {
	service_name = data.ovh_cloud_project_capabilities_containerregistry_filter.regcap.service_name
    plan_id      = data.ovh_cloud_project_capabilities_containerregistry_filter.regcap.id
	name         = "%s"
    region       = data.ovh_cloud_project_capabilities_containerregistry_filter.regcap.region
}

data "ovh_cloud_project_containerregistry" "reg" {
  service_name = ovh_cloud_project_containerregistry.reg.service_name
  registry_id  = ovh_cloud_project_containerregistry.reg.id
}
`

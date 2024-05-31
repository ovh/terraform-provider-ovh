package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCloudProjectContainerRegistryIPRestrictionsManagementDataSourceConfig = `
data "ovh_cloud_project_capabilities_containerregistry_filter" "registryCap" {
  service_name = "%s"
  plan_name    = "SMALL"
  region       = "%s"
}

resource "ovh_cloud_project_containerregistry" "registry" {
  service_name = data.ovh_cloud_project_capabilities_containerregistry_filter.registryCap.service_name
  plan_id      = data.ovh_cloud_project_capabilities_containerregistry_filter.registryCap.id
  name         = "%s"
  region       = data.ovh_cloud_project_capabilities_containerregistry_filter.registryCap.region
}

resource "ovh_cloud_project_containerregistry_ip_restrictions_management" "my-mgt-iprestrictions" {
  service_name = ovh_cloud_project_containerregistry.registry.service_name
  registry_id  = ovh_cloud_project_containerregistry.registry.id
	
  ip_restrictions = [
    {
      ip_block = "121.121.121.121/32"
      description = "my awesome ip"
    }
  ]
  depends_on = [
    ovh_cloud_project_containerregistry.registry
  ]
}

data "ovh_cloud_project_containerregistry_ip_restrictions_management" "mgt-iprestrictions-data" {
  service_name = ovh_cloud_project_containerregistry.registry.service_name
  registry_id  = ovh_cloud_project_containerregistry.registry.id

  depends_on = [
    ovh_cloud_project_containerregistry_ip_restrictions_management.my-mgt-iprestrictions
  ]
}
`

func TestAccCloudProjectContainerIPRestrictionsManagementDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_CONTAINERREGISTRY_REGION_TEST")
	registryName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccCloudProjectContainerRegistryIPRestrictionsManagementDataSourceConfig,
		serviceName,
		region,
		registryName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckContainerRegistry(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_containerregistry_ip_restrictions_management.mgt-iprestrictions-data", "ip_restrictions.0.ip_block", "121.121.121.121/32"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_containerregistry_ip_restrictions_management.mgt-iprestrictions-data", "ip_restrictions.0.description", "my awesome ip"),
				),
			},
		},
	})
}

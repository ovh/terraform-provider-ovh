package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectContainerRegistryIAMDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_CONTAINERREGISTRY_REGION_TEST")
	registryName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccCloudProjectContainerRegistryIAMDataSourceConfig,
		serviceName,
		region,
		registryName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckContainerRegistryIAM(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_containerregistry_iam.iamData", "iam_enabled", "true"),
				),
			},
		},
	})
}

var testAccCloudProjectContainerRegistryIAMDataSourceConfig = `
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

resource "ovh_cloud_project_containerregistry_iam" "iam" {
	service_name = ovh_cloud_project_containerregistry.registry.service_name
	registry_id  = ovh_cloud_project_containerregistry.registry.id
	delete_users = "true"

	depends_on = [
		ovh_cloud_project_containerregistry.registry
	]
}

data "ovh_cloud_project_containerregistry_iam" "iamData" {
    service_name = ovh_cloud_project_containerregistry.registry.service_name
    registry_id = ovh_cloud_project_containerregistry.registry.id

	depends_on = [
		ovh_cloud_project_containerregistry_iam.iam
	]
}
`

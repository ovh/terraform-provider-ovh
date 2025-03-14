package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var testAccCloudProjectContainerRegistryIAMConfig = `
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

resource "ovh_cloud_project_containerregistry_iam" "my-iam" {
	service_name = ovh_cloud_project_containerregistry.registry.service_name
	registry_id  = ovh_cloud_project_containerregistry.registry.id

	delete_users = true
}
`

func TestAccCloudProjectContainerRegistryIAM_full(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	registryName := acctest.RandomWithPrefix(test_prefix)
	region := os.Getenv("OVH_CLOUD_PROJECT_CONTAINERREGISTRY_REGION_TEST")
	resourceName := "ovh_cloud_project_containerregistry_iam.my-iam"

	config := fmt.Sprintf(
		testAccCloudProjectContainerRegistryIAMConfig,
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
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_iam.my-iam", "delete_users", "true"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_containerregistry_iam.my-iam", "iam_enabled"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_users"},
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s", serviceName, state.RootModule().Resources[resourceName].Primary.Attributes["registry_id"]), nil
				},
			},
			{
				Config:       config,
				Destroy:      true,
				ResourceName: resourceName,
			},
		},
	})
}

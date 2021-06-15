package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudProjectContainerRegistryUser_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	regName := acctest.RandomWithPrefix(test_prefix)
	region := "GRA"
	config := fmt.Sprintf(
		testAccCloudProjectContainerRegistryUserConfig,
		serviceName,
		region,
		regName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_containerregistry.reg", "name", regName),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_containerregistry.reg", "region", region),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_containerregistry_user.user",
						"password",
					),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_containerregistry_user.user",
						"id",
					),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_containerregistry_user.user",
						"user", "foobar",
					),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_containerregistry_user.user",
						"email", "foo@bar.com",
					),
				),
			},
		},
	})
}

const testAccCloudProjectContainerRegistryUserConfig = `
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

resource "ovh_cloud_project_containerregistry_user" "user" {
	service_name = ovh_cloud_project_containerregistry.reg.service_name
    registry_id  = ovh_cloud_project_containerregistry.reg.id
	email        = "foo@bar.com"
    login        = "foobar"
}
`

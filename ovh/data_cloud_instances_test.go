package ovh

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudInstances_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_INSTANCE_REGION_TEST")
	flavorID := os.Getenv("OVH_INSTANCE_FLAVOR_ID_TEST")
	imageID := os.Getenv("OVH_INSTANCE_IMAGE_ID_TEST")
	name := acctest.RandomWithPrefix("test-inst-list")

	config := testAccCloudInstanceConfig(serviceName, region, flavorID, imageID, name) + `
data "ovh_cloud_instances" "all" {
  service_name = ovh_cloud_instance.test.service_name
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstance(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_cloud_instances.all", "instances.#"),
					// At least the instance created above must be listed.
					resource.TestCheckResourceAttrWith("data.ovh_cloud_instances.all", "instances.#", testAccCheckCloudPublicIPListNotEmpty),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_instances.all", "instances.0.id"),
				),
			},
		},
	})
}

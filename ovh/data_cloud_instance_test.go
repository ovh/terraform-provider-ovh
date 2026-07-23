package ovh

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudInstance_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	name := acctest.RandomWithPrefix("test-inst-ds")

	config := testAccCloudInstanceConfig(serviceName, region, flavorID, imageID, name) + `
data "ovh_cloud_instance" "test" {
  service_name = ovh_cloud_instance.test.service_name
  id           = ovh_cloud_instance.test.id
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstance(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_instance.test", "name", name),
					resource.TestCheckResourceAttr("data.ovh_cloud_instance.test", "flavor_id", flavorID),
					resource.TestCheckResourceAttr("data.ovh_cloud_instance.test", "resource_status", "READY"),
					// Observed state is fully populated on a single-instance read.
					resource.TestCheckResourceAttrSet("data.ovh_cloud_instance.test", "current_state.power_state"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_instance.test", "current_state.flavor.id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_instance.test", "current_state.location.region"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_instance.test", "current_state.networks.0.id"),
					// The data source must resolve to the same instance as the resource.
					resource.TestCheckResourceAttrPair("data.ovh_cloud_instance.test", "id", "ovh_cloud_instance.test", "id"),
					resource.TestCheckResourceAttrPair("data.ovh_cloud_instance.test", "current_state.flavor.id", "ovh_cloud_instance.test", "current_state.flavor.id"),
				),
			},
		},
	})
}

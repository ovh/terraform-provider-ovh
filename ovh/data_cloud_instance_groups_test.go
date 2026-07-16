package ovh

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudInstanceGroups_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_INSTANCE_REGION_TEST")
	name := acctest.RandomWithPrefix("test-grp-list")

	config := testAccCloudInstanceGroupConfig(serviceName, region, name) + `
data "ovh_cloud_instance_groups" "all" {
  service_name = ovh_cloud_instance_group.test.service_name
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceGroup(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_cloud_instance_groups.all", "instance_groups.#"),
					// At least the group created above must be listed.
					resource.TestCheckResourceAttrWith("data.ovh_cloud_instance_groups.all", "instance_groups.#", testAccCheckCloudPublicIPListNotEmpty),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_instance_groups.all", "instance_groups.0.id"),
				),
			},
		},
	})
}

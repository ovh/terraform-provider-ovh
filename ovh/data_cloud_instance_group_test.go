package ovh

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudInstanceGroup_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	name := acctest.RandomWithPrefix("test-grp-ds")

	config := testAccCloudInstanceGroupConfig(serviceName, region, name) + `
data "ovh_cloud_instance_group" "test" {
  service_name = ovh_cloud_instance_group.test.service_name
  id           = ovh_cloud_instance_group.test.id
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceGroup(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_instance_group.test", "name", name),
					resource.TestCheckResourceAttr("data.ovh_cloud_instance_group.test", "region", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_instance_group.test", "policy", "ANTI_AFFINITY"),
					// The data source must resolve to the same group as the resource.
					resource.TestCheckResourceAttrPair("data.ovh_cloud_instance_group.test", "id", "ovh_cloud_instance_group.test", "id"),
				),
			},
		},
	})
}

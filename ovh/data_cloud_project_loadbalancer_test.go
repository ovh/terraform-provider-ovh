package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudProjectBalancer_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudRegionLoadbalancer(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_cloud_project_loadbalancer" "lb" {
						service_name = "%s"
						region_name  = "%s"
						id           = "%s"
					}
				`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"), os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST"), os.Getenv("OVH_CLOUD_PROJECT_LOADBALANCER_TEST")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_project_loadbalancer.lb", "region_name", os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_loadbalancer.lb", "name"),
				),
			},
		},
	})
}

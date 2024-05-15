package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudProjectBalancers_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudRegion(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_cloud_project_loadbalancers" "lbs" {
						service_name = "%s"
						region_name  = "%s"
					}
				`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"), os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_loadbalancers.lbs", "loadbalancers.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_loadbalancers.lbs", "loadbalancers.0.vip_address"),
				),
			},
		},
	})
}

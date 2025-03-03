package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudProjectRancherPlan_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_cloud_project_rancher_plan" "plans" {
						project_id = "%s"
					}
				`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_project_rancher_plan.plans", "project_id", os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_rancher_plan.plans", "plans.#"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_rancher_plan.plans", "plans.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_rancher_plan.plans", "plans.0.status"),
				),
			},
		},
	})
}

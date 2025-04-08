package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudProjectRancher_basic(t *testing.T) {
	projectID := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	rancherID := os.Getenv("OVH_CLOUD_PROJECT_RANCHER_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckRancher(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_cloud_project_rancher" "rancher" {
						project_id = "%s"
						id         = "%s"
					}
				`, projectID, rancherID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_project_rancher.rancher", "project_id", projectID),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_rancher.rancher", "id", rancherID),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_rancher.rancher", "current_state.region", "EU_WEST_SBG"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_rancher.rancher", "resource_status", "READY"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_rancher.rancher", "target_spec.plan", "OVHCLOUD_EDITION"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_rancher.rancher", "current_state.version"),
				),
			},
		},
	})
}

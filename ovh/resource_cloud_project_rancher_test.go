package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectRancher_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "ovh_cloud_project_rancher" "ranch" {
					project_id = "%s"
					target_spec = {
						name = "MyFirstRancher"
						plan = "STANDARD"
					}
				}
				`, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_rancher.ranch", "project_id", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_rancher.ranch", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_cloud_project_rancher.ranch", "current_state.name", "MyFirstRancher"),
					resource.TestCheckResourceAttr("ovh_cloud_project_rancher.ranch", "current_state.plan", "STANDARD"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_rancher.ranch", "current_state.version"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_rancher.ranch", "current_state.region"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_rancher.ranch", "current_state.url"),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "ovh_cloud_project_rancher" "ranch" {
					project_id = "%s"
					target_spec = {
						name = "MyFirstRancherUpdated"
						plan = "OVHCLOUD_EDITION"
					}
				}
				`, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_rancher.ranch", "project_id", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_rancher.ranch", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_cloud_project_rancher.ranch", "current_state.name", "MyFirstRancherUpdated"),
					resource.TestCheckResourceAttr("ovh_cloud_project_rancher.ranch", "current_state.plan", "OVHCLOUD_EDITION"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_rancher.ranch", "current_state.version"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_rancher.ranch", "current_state.region"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_rancher.ranch", "current_state.url"),
				),
			},
		},
	})
}

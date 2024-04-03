package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudProject_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "ovh_cloud_project" "project" {
					service_name = "%s"
				}
				`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_project.project", "access", "full"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project.project", "plan_code", "project.2018"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project.project", "status", "ok"),
				),
			},
		},
	})
}

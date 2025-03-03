package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudProjectRancherVersion_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_cloud_project_rancher_version" "versions" {
						project_id = "%s"
					}
				`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_project_rancher_version.versions", "project_id", os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_rancher_version.versions", "versions.#"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_rancher_version.versions", "versions.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_rancher_version.versions", "versions.0.status"),
				),
			},
		},
	})
}

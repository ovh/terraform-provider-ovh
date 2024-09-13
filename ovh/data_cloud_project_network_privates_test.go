package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudProjectNetworkPrivates_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudNetworkPrivate(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_cloud_project_network_privates" "private" {
						service_name = "%s"
					}
				`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_project_network_privates.private", "service_name", os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_network_privates.private", "networks.#"),
				),
			},
		},
	})
}

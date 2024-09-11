package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudProjectNetworkPrivate_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudNetworkPrivate(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_cloud_project_network_private" "private" {
						service_name = "%s"
						network_id = "%s"
					}
				`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"), os.Getenv("OVH_CLOUD_PROJECT_PRIVATE_NETWORK_TEST")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_project_network_private.private", "service_name", os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_network_private.private", "network_id", os.Getenv("OVH_CLOUD_PROJECT_PRIVATE_NETWORK_TEST")),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_network_private.private", "name"),
				),
			},
		},
	})
}

package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testCreateNetworkPrivate = fmt.Sprintf(`
resource "ovh_cloud_project_network_private" "testnetwork" {
      service_name = "%s"
	  name         = "network_test"
	  regions      = ["GRA11", "GRA9"]
	  vlan_id      = 0
	}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

var testUpdateNetworkPrivate = fmt.Sprintf(`
resource "ovh_cloud_project_network_private" "testnetwork" {
      service_name = "%s"
	  name         = "network_test"
	  regions      = ["GRA11"]
	  vlan_id      = 0
	}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

func TestAccCloudProjectNetworkPrivate_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCreateNetworkPrivate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_network_private.testnetwork", "regions.#", "2"),
				),
			},
			{
				Config: testUpdateNetworkPrivate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_network_private.testnetwork", "regions.#", "1"),
				),
			},
		},
	})
}

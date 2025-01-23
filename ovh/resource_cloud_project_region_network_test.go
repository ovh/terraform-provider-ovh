package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectRegionNetwork_basic(t *testing.T) {
	config := fmt.Sprintf(`
	resource "ovh_cloud_project_region_network" "net" {
		service_name = "%s"
		region_name  = "EU-SOUTH-LZ-MAD-A"
		name         = "MadriNet"
		subnet       = {
			cidr              = "10.0.2.0/24"
			enable_dhcp       = true
			enable_gateway_ip = false
			ip_version        = 4
		}
	}
	`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_region_network.net", "name", "MadriNet"),
					resource.TestCheckResourceAttr("ovh_cloud_project_region_network.net", "visibility", "private"),
					resource.TestCheckResourceAttr("ovh_cloud_project_region_network.net", "subnet.cidr", "10.0.2.0/24"),
				),
			},
		},
	})
}

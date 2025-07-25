package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectGatewayDataSource_basic(t *testing.T) {
	gatewayName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
			resource "ovh_cloud_project_network_private" "mypriv" {
				service_name  = "%s"
				vlan_id       = "%d"
				name          = "%s"
				regions       = ["GRA11"]
			}

			resource "ovh_cloud_project_network_private_subnet" "myprivsub" {
				service_name  = ovh_cloud_project_network_private.mypriv.service_name
				network_id    = ovh_cloud_project_network_private.mypriv.id
				region        = "GRA11"
				start         = "10.0.0.2"
				end           = "10.0.0.8"
				network       = "10.0.0.0/24"
				dhcp          = true
			}

			resource "ovh_cloud_project_gateway" "gateway" {
				service_name = ovh_cloud_project_network_private.mypriv.service_name
				name          = "%s"
				model         = "s"
				region        = ovh_cloud_project_network_private_subnet.myprivsub.region
				network_id    = tolist(ovh_cloud_project_network_private.mypriv.regions_attributes[*].openstackid)[0]
				subnet_id     = ovh_cloud_project_network_private_subnet.myprivsub.id
			}

			data "ovh_cloud_project_gateway" "gateway" {
				service_name = ovh_cloud_project_gateway.gateway.service_name
				region       = ovh_cloud_project_gateway.gateway.region
				id           = ovh_cloud_project_gateway.gateway.id
			}
		`,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		acctest.RandIntRange(100, 200),
		acctest.RandomWithPrefix(test_prefix),
		gatewayName,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_project_gateway.gateway", "name", gatewayName),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_gateway.gateway", "external_information.ips.#"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_gateway.gateway", "external_information.ips.0.ip"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_gateway.gateway", "external_information.ips.0.subnet_id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_gateway.gateway", "interfaces.#"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_gateway.gateway", "interfaces.0.id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_gateway.gateway", "interfaces.0.ip"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_gateway.gateway", "model"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_gateway.gateway", "status"),
				),
			},
		},
	})
}

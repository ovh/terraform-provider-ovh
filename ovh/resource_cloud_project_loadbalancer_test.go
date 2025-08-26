package ovh

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectLoadBalancer_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	vlanId := rand.Intn(4000)

	config := fmt.Sprintf(`
		resource "ovh_cloud_project_network_private" "mypriv" {
			service_name = "%s"
			name         = "network_test"
			regions      = ["GRA11", "GRA9"]
			vlan_id      = %d
		}

		resource "ovh_cloud_project_network_private_subnet" "myprivsub" {
			service_name  = ovh_cloud_project_network_private.mypriv.service_name
			network_id    = ovh_cloud_project_network_private.mypriv.id
			region        = "GRA9"
			start         = "10.0.0.2"
			end           = "10.0.255.254"
			network       = "10.0.0.0/16"
			dhcp          = true
		  }

		resource "ovh_cloud_project_loadbalancer" "lb" {
			service_name = ovh_cloud_project_network_private_subnet.myprivsub.service_name
			region_name = "GRA9"
			flavor_id = "2d4bc92c-38fc-4b50-9484-8351ab0c4e69"
			network = {
			private = {
				network = {
					id = element([for region in ovh_cloud_project_network_private.mypriv.regions_attributes: region if "${region.region}" == "GRA9"], 0).openstackid
					subnet_id = ovh_cloud_project_network_private_subnet.myprivsub.id
				}
			}
			}
			description = "A great load balancer"
			listeners = [
				{
					port = "34568"
					protocol = "tcp"
					pool = {
						algorithm      = "roundRobin"
						name           = "TestPool"
						protocol       = "http"
						members = [
							{
								name          = "web"
								address       = "1.2.3.4"
								protocol_port = 80
							}
						]
					}
				}
			]
		}
	`, serviceName, vlanId)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_loadbalancer.lb", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_loadbalancer.lb", "description", "A great load balancer"),
					resource.TestCheckResourceAttr("ovh_cloud_project_loadbalancer.lb", "operating_status", "online"),
					resource.TestCheckResourceAttr("ovh_cloud_project_loadbalancer.lb", "region", "GRA9"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_loadbalancer.lb", "vip_network_id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_loadbalancer.lb", "vip_subnet_id"),
				),
			},
		},
	})
}

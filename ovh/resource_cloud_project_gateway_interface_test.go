package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccCloudProjectGatewayInterfaceConfig = `
data "ovh_me" "myaccount" {}

data "ovh_order_cart" "mycart" {
	ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
}

data "ovh_order_cart_product_plan" "cloud" {
	cart_id        = data.ovh_order_cart.mycart.id
	price_capacity = "renew"
	product        = "cloud"
	plan_code      = "project.2018"
}

resource "ovh_cloud_project" "cloud" {
	ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
	description    = "Cloud project for gateway interface test"

	plan {
		duration     = data.ovh_order_cart_product_plan.cloud.selected_price.0.duration
		plan_code    = data.ovh_order_cart_product_plan.cloud.plan_code
		pricing_mode = data.ovh_order_cart_product_plan.cloud.selected_price.0.pricing_mode
	}
}

resource "ovh_vrack_cloudproject" "attach" {
	service_name = "%s"
	project_id   = ovh_cloud_project.cloud.project_id
}

resource "ovh_cloud_project_network_private" "mypriv" {
	service_name  = ovh_vrack_cloudproject.attach.project_id
	vlan_id       = "%d"
	name          = "%s"
	regions       = ["GRA9"]
}

resource "ovh_cloud_project_network_private_subnet" "myprivsub" {
	service_name  = ovh_cloud_project_network_private.mypriv.service_name
	network_id    = ovh_cloud_project_network_private.mypriv.id
	region        = "GRA9"
	start         = "10.0.0.2"
	end           = "10.0.0.8"
	network       = "10.0.0.0/24"
	dhcp          = true
}

resource "ovh_cloud_project_network_private_subnet" "my_other_privsub" {
	service_name  = ovh_cloud_project_network_private.mypriv.service_name
	network_id    = ovh_cloud_project_network_private.mypriv.id
	region        = "GRA9"
	start         = "10.0.1.10"
	end           = "10.0.1.254"
	network       = "10.0.1.0/24"
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

resource "ovh_cloud_project_gateway_interface" "interface" {
	service_name = ovh_cloud_project_network_private.mypriv.service_name
	region       = ovh_cloud_project_network_private_subnet.my_other_privsub.region
	id           = ovh_cloud_project_gateway.gateway.id
	subnet_id    = ovh_cloud_project_network_private_subnet.my_other_privsub.id
}
`

func TestAccCloudProjectGatewayInterface_basic(t *testing.T) {
	config := fmt.Sprintf(
		testAccCloudProjectGatewayInterfaceConfig,
		os.Getenv("OVH_VRACK_SERVICE_TEST"),
		acctest.RandIntRange(100, 200),
		acctest.RandomWithPrefix(test_prefix),
		acctest.RandomWithPrefix(test_prefix),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_gateway_interface.interface", "ip"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_gateway_interface.interface", "subnet_id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_gateway_interface.interface", "network_id"),
				),
			},
			{
				Config:            config,
				ImportStateVerify: true,
			},
		},
	})
}

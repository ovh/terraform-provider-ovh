package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccOvhcloudConnectPopConfig = `
data "ovh_ovhcloud_connect" "occ" {
	service_name = "%s"
}

resource "ovh_ovhcloud_connect_pop_config" "pop" {
	service_name = data.ovh_ovhcloud_connect.occ.service_name
	interface_id = tolist(data.ovh_ovhcloud_connect.occ.interface_list)[0]
	type = "l3"
	customer_bgp_area = 65005
	ovh_bgp_area = 65006
	subnet = "10.10.12.0/30"
}
`

func TestAccOvhcloudConnectPopConfig_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_OCC_SERVICE_TEST")

	config := fmt.Sprintf(testAccOvhcloudConnectPopConfig, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			checkEnvOrSkip(t, "OVH_OCC_SERVICE_TEST")
			testAccPreCheckCredentials(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_ovhcloud_connect_pop_config.pop", "type", "l3"),
					resource.TestCheckResourceAttr("ovh_ovhcloud_connect_pop_config.pop", "customer_bgp_area", "65005"),
					resource.TestCheckResourceAttr("ovh_ovhcloud_connect_pop_config.pop", "ovh_bgp_area", "65006"),
					resource.TestCheckResourceAttr("ovh_ovhcloud_connect_pop_config.pop", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_ovhcloud_connect_pop_config.pop", "subnet", "10.10.12.0/30"),
					resource.TestCheckResourceAttr("ovh_ovhcloud_connect_pop_config.pop", "status", "active"),
				),
			},
			{
				Config:  config,
				Destroy: true,
			},
		},
	})
}

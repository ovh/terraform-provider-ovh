package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccOvhcloudConnectPopDatacenterConfig = `
variable "service_name" {
    type        = string
    default     = "%s"
}

data "ovh_ovhcloud_connect_config_pops" "pop_cfgs" {
    service_name = var.service_name
}

resource "ovh_ovhcloud_connect_pop_datacenter_config" "dc" {
    service_name = var.service_name
    config_pop_id = tolist(data.ovh_ovhcloud_connect_config_pops.pop_cfgs.pop_configs)[0].id
    datacenter_id = 6
    ovh_bgp_area = 65008
    subnet = "10.0.0.0/28"
}
`

func TestAccOvhcloudConnectPopDatacenterConfig_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_OCC_SERVICE_TEST")

	config := fmt.Sprintf(testAccOvhcloudConnectPopDatacenterConfig, serviceName)

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
					resource.TestCheckResourceAttr("ovh_ovhcloud_connect_pop_datacenter_config.dc", "ovh_bgp_area", "65008"),
					resource.TestCheckResourceAttr("ovh_ovhcloud_connect_pop_datacenter_config.dc", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_ovhcloud_connect_pop_datacenter_config.dc", "subnet", "10.0.0.0/28"),
					resource.TestCheckResourceAttr("ovh_ovhcloud_connect_pop_datacenter_config.dc", "status", "active"),
				),
			},
			{
				Config:  config,
				Destroy: true,
			},
		},
	})
}

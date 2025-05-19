package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccOvhcloudConnectPopDatacenterExtraConfig = `
variable "service_name" {
    type        = string
    default     = "%s"
}

data "ovh_ovhcloud_connect_config_pops" "pop_cfgs" {
    service_name = var.service_name
}

data "ovh_ovhcloud_connect_config_pop_datacenters" "datacenter_cfgs" {
  service_name = var.service_name
  config_pop_id = tolist(data.ovh_ovhcloud_connect_config_pops.pop_cfgs.pop_configs)[0].id
}

resource "ovh_ovhcloud_connect_pop_datacenter_extra_config" "extra" {
    service_name = var.service_name
    config_pop_id = tolist(data.ovh_ovhcloud_connect_config_pops.pop_cfgs.pop_configs)[0].id
    config_datacenter_id = tolist(data.ovh_ovhcloud_connect_config_pop_datacenters.datacenter_cfgs.datacenter_configs)[0].id
    type = "network"
    next_hop = "10.0.0.5"
    subnet = "192.168.2.0/24"
}
`

func TestAccOvhcloudConnectPopDatacenterExtraConfig_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_OCC_SERVICE_TEST")

	config := fmt.Sprintf(testAccOvhcloudConnectPopDatacenterExtraConfig, serviceName)

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
					resource.TestCheckResourceAttr("ovh_ovhcloud_connect_pop_datacenter_extra_config.extra", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_ovhcloud_connect_pop_datacenter_extra_config.extra", "next_hop", "10.0.0.5"),
					resource.TestCheckResourceAttr("ovh_ovhcloud_connect_pop_datacenter_extra_config.extra", "subnet", "192.168.2.0/24"),
					resource.TestCheckResourceAttr("ovh_ovhcloud_connect_pop_datacenter_extra_config.extra", "status", "active"),
				),
			},
			{
				Config:  config,
				Destroy: true,
			},
		},
	})
}

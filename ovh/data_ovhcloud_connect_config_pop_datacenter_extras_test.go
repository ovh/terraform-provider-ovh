package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceOvhCloudConnectConfigPopDatacenterExtras_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_OCC_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			checkEnvOrSkip(t, "OVH_OCC_SERVICE_TEST")
			testAccPreCheckCredentials(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_ovhcloud_connect_config_pops" "pop_cfgs" {
						service_name = "%s"
					}
					data "ovh_ovhcloud_connect_config_pop_datacenters" "datacenter_cfgs" {
						service_name = data.ovh_ovhcloud_connect_config_pops.pop_cfgs.service_name
						config_pop_id = tolist(data.ovh_ovhcloud_connect_config_pops.pop_cfgs.pop_configs)[0].id
					}
					data "ovh_ovhcloud_connect_config_pop_datacenter_extras" "extra_cfgs" {
						service_name = data.ovh_ovhcloud_connect_config_pops.pop_cfgs.service_name
						config_pop_id = tolist(data.ovh_ovhcloud_connect_config_pops.pop_cfgs.pop_configs)[0].id
						config_datacenter_id = tolist(data.ovh_ovhcloud_connect_config_pop_datacenters.datacenter_cfgs.datacenter_configs)[0].id
					}
				`, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_ovhcloud_connect_config_pop_datacenter_extras.extra_cfgs", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect_config_pop_datacenter_extras.extra_cfgs", "config_pop_id"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect_config_pop_datacenter_extras.extra_cfgs", "config_datacenter_id"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect_config_pop_datacenter_extras.extra_cfgs", "extra_configs.0.id"),
					resource.TestCheckResourceAttr("data.ovh_ovhcloud_connect_config_pop_datacenter_extras.extra_cfgs", "extra_configs.0.service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect_config_pop_datacenter_extras.extra_cfgs", "extra_configs.0.status"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect_config_pop_datacenter_extras.extra_cfgs", "extra_configs.0.type"),
				),
			},
		},
	})
}

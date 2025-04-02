package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceOvhCloudConnectConfigPop_basic(t *testing.T) {
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
					data "ovh_ovhcloud_connect_config_pops" "pop_configs" {
						service_name = "%s"
					}
				`, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_ovhcloud_connect_config_pops.pop_configs", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect_config_pops.pop_configs", "pop_configs.0.id"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect_config_pops.pop_configs", "pop_configs.0.interface_id"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect_config_pops.pop_configs", "pop_configs.0.pop_id"),
					resource.TestCheckResourceAttr("data.ovh_ovhcloud_connect_config_pops.pop_configs", "pop_configs.0.service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect_config_pops.pop_configs", "pop_configs.0.status"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect_config_pops.pop_configs", "pop_configs.0.type"),
				),
			},
		},
	})
}

package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceOvhCloudConnectDatacenters_basic(t *testing.T) {
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
					data "ovh_ovhcloud_connect_datacenters" "dcs" {
						service_name = "%s"
					}
				`, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_ovhcloud_connect_datacenters.dcs", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect_datacenters.dcs", "datacenters.0.available"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect_datacenters.dcs", "datacenters.0.id"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect_datacenters.dcs", "datacenters.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect_datacenters.dcs", "datacenters.0.region"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect_datacenters.dcs", "datacenters.0.region_type"),
				),
			},
		},
	})
}

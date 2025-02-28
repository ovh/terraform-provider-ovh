package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceOvhCloudConnect_basic(t *testing.T) {
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
					data "ovh_ovhcloud_connect" "occ" {
						service_name = "%s"
					}
				`, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_ovhcloud_connect.occ", "uuid", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect.occ", "bandwidth"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect.occ", "description"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect.occ", "interface_list.0"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect.occ", "pop"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect.occ", "port_quantity"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect.occ", "product"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect.occ", "provider_name"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connect.occ", "status"),
				),
			},
		},
	})
}

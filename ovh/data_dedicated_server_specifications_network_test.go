package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDedicatedServerSpecificationsNetworkDataSource_basic(t *testing.T) {
	testAccDedicatedServerSpecificationsNetworkDatasourceConfig_Basic := fmt.Sprintf(`
	data "ovh_dedicated_server_specifications_network" "spec" {
		service_name = "%s"
	}`, os.Getenv("OVH_DEDICATED_SERVER"))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCredentials(t)
			testAccPreCheckDedicatedServer(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerSpecificationsNetworkDatasourceConfig_Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_dedicated_server_specifications_network.spec", "bandwidth.internet_to_ovh.value", "500"),
					resource.TestCheckResourceAttr("data.ovh_dedicated_server_specifications_network.spec", "bandwidth.internet_to_ovh.unit", "Mbps"),
					resource.TestCheckResourceAttr("data.ovh_dedicated_server_specifications_network.spec", "switching.name", "bhs6-sdtor46a-n3"),
				),
			},
		},
	})
}

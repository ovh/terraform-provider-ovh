package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDedicatedServerOrderableBandwidthDataSource_basic(t *testing.T) {
	testAccDedicatedServerOrderableBandwidthDatasourceConfig_Basic := fmt.Sprintf(`
	data "ovh_dedicated_server_orderable_bandwidth" "bp" {
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
				Config: testAccDedicatedServerOrderableBandwidthDatasourceConfig_Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_dedicated_server_orderable_bandwidth.bp", "orderable"),
					resource.TestCheckResourceAttrSet("data.ovh_dedicated_server_orderable_bandwidth.bp", "platinium"),
					resource.TestCheckResourceAttrSet("data.ovh_dedicated_server_orderable_bandwidth.bp", "premium"),
					resource.TestCheckResourceAttrSet("data.ovh_dedicated_server_orderable_bandwidth.bp", "ultimate"),
				),
			},
		},
	})
}

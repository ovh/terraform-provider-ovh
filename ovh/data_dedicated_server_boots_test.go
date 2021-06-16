package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDedicatedServerBootsDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCredentials(t)
			testAccPreCheckDedicatedServer(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerBootsDatasourceConfig(""),
				Check:  resource.TestCheckOutput("test", "true"),
			},
			{
				Config: testAccDedicatedServerBootsDatasourceConfig("harddisk"),
				Check:  resource.TestCheckOutput("test", "true"),
			},
			{
				Config: testAccDedicatedServerBootsDatasourceConfig("rescue"),
				Check:  resource.TestCheckOutput("test", "true"),
			},
			{
				Config: testAccDedicatedServerBootsDatasourceConfig("network"),
				Check:  resource.TestCheckOutput("test", "true"),
			},
		},
	})
}

func testAccDedicatedServerBootsDatasourceConfig(filter string) string {
	dedicated_server := os.Getenv("OVH_DEDICATED_SERVER")
	if filter == "" {
		return fmt.Sprintf(
			testAccDedicatedServerBootsDatasourceConfig_Basic,
			dedicated_server,
		)
	}
	return fmt.Sprintf(
		testAccDedicatedServerBootsDatasourceConfig_Filter,
		dedicated_server,
		filter,
	)

}

const testAccDedicatedServerBootsDatasourceConfig_Basic = `
data "ovh_dedicated_server_boots" "boots" {
  service_name = "%s"
}

output test { value = tostring(length(data.ovh_dedicated_server_boots.boots.result) > 0 )}
`
const testAccDedicatedServerBootsDatasourceConfig_Filter = `
data "ovh_dedicated_server_boots" "boots" {
  service_name = "%s"
  boot_type    = "%s"
}

output test { value = tostring(length(data.ovh_dedicated_server_boots.boots.result) > 0 )}
`

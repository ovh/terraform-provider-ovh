package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDedicatedServerDataSource_basic(t *testing.T) {
	dedicated_server := os.Getenv("OVH_DEDICATED_SERVER_SERVICE_NAME")
	config := fmt.Sprintf(testAccDedicatedServerDatasourceConfig_Basic, dedicated_server)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDedicatedServer(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_dedicated_server.server", "name", dedicated_server),
					resource.TestCheckResourceAttr(
						"data.ovh_dedicated_server.server", "service_name", dedicated_server),
				),
			},
		},
	})
}

const testAccDedicatedServerDatasourceConfig_Basic = `
data "ovh_dedicated_server" "server" {
  service_name  = "%s"
}
`

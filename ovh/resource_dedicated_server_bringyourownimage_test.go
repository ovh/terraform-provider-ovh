package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDedicatedServerBringYourOwnImage(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCredentials(t)
			testAccPreCheckDedicatedServer(t)
			testAccPreCheckByoiImageUrl(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerBringYourOwnImageConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server.server", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server.server", "monitoring", "true"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server.server", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_bringyourownimage.os_install", "function", "reinstallServer"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_bringyourownimage.os_install", "status", "done"),
				),
			},
		},
	})
}

func testAccDedicatedServerBringYourOwnImageConfig() string {
	dedicated_server := os.Getenv("OVH_DEDICATED_SERVER")
	url := os.Getenv("OVH_BYOI_IMAGE_URL")
	hostName := acctest.RandomWithPrefix(test_prefix)

	return fmt.Sprintf(`
	data ovh_dedicated_server "server" {
		service_name = "%s"
		boot_type    = "harddisk"
	}
	
	resource ovh_dedicated_server_bringyourownimage "os_install" {
		service_name = data.ovh_dedicated_server.server.name
		url = "%s"
		type = "qcow2"
		config_drive {
			enable = true
			hostname = "%s"
		}
	}
	`,
		dedicated_server,
		url,
		hostName,
	)
}

func testAccPreCheckByoiImageUrl(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_BYOI_IMAGE_URL")
}

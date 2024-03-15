package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDedicatedServerSpecificationsHardwareDataSource_basic(t *testing.T) {
	testAccDedicatedServerSpecificationsHardwareDatasourceConfig_Basic := fmt.Sprintf(`
	data "ovh_dedicated_server_specifications_hardware" "spec" {
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
				Config: testAccDedicatedServerSpecificationsHardwareDatasourceConfig_Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_dedicated_server_specifications_hardware.spec", "number_of_processors", "1"),
					resource.TestCheckResourceAttr("data.ovh_dedicated_server_specifications_hardware.spec", "motherboard", "E3C246D4U2-2T"),
					resource.TestCheckResourceAttr("data.ovh_dedicated_server_specifications_hardware.spec", "description", "RISE-1 - Intel Xeon-E 2136"),
					resource.TestCheckResourceAttr("data.ovh_dedicated_server_specifications_hardware.spec", "processor_name", "XeonE-2236"),
					resource.TestCheckResourceAttr("data.ovh_dedicated_server_specifications_hardware.spec", "memory_size.value", "32768"),
					resource.TestCheckResourceAttr("data.ovh_dedicated_server_specifications_hardware.spec", "memory_size.unit", "MB"),
					resource.TestCheckResourceAttr("data.ovh_dedicated_server_specifications_hardware.spec", "disk_groups.0.disk_type", "NVME"),
					resource.TestCheckResourceAttr("data.ovh_dedicated_server_specifications_hardware.spec", "disk_groups.0.disk_size.value", "512"),
					resource.TestCheckResourceAttr("data.ovh_dedicated_server_specifications_hardware.spec", "disk_groups.0.disk_size.unit", "GB"),
				),
			},
		},
	})
}

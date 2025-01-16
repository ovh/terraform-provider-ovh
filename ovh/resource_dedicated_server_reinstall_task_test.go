package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDedicatedServerReinstall_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCredentials(t)
			testAccPreCheckDedicatedServer(t)
		},
		Providers: testAccProviders,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				VersionConstraint: "0.10.0",
				Source:            "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerReinstallConfig("basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "monitoring", "true"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_reinstall_task.server_reinstall", "function", "reinstallServer"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_reinstall_task.server_reinstall", "status", "done"),
				),
			},
		},
	})
}

func TestAccDedicatedServerReinstall_rebootondestroy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCredentials(t)
			testAccPreCheckDedicatedServer(t)
		},
		Providers: testAccProviders,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				VersionConstraint: "0.10.0",
				Source:            "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerReinstallConfig("rebootondestroy"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "monitoring", "true"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_reinstall_task.server_reinstall", "function", "reinstallServer"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_reinstall_task.server_reinstall", "status", "done"),
				),
			},
		},
	})
}

func testAccDedicatedServerReinstallConfig(config string) string {
	dedicated_server := os.Getenv("OVH_DEDICATED_SERVER")
	sshKey := os.Getenv("OVH_SSH_KEY")
	if sshKey == "" {
		sshKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIrODOo0SvY5f0TlQNvGHIRKzr4bHPa+D5bYF18RiOgP email@example.com"
	}

	if config == "rebootondestroy" {
		return fmt.Sprintf(
			testAccDedicatedServerReinstallConfig_RebootOnDestroy,
			dedicated_server,
		)
	}

	return fmt.Sprintf(
		testAccDedicatedServerReinstallConfig_Basic,
		dedicated_server,
	)

}

const testAccDedicatedServerReinstallConfig_Basic = `
data ovh_dedicated_server_boots "harddisk" {
  service_name = "%s"
  boot_type    = "harddisk"
}

resource ovh_dedicated_server_update "server" {
  service_name = data.ovh_dedicated_server_boots.harddisk.service_name
  boot_id      = data.ovh_dedicated_server_boots.harddisk.result[0]
  monitoring   = true
  state        = "ok"
}

resource ovh_dedicated_server_reinstall_task "server_reinstall" {
  service_name = data.ovh_dedicated_server_boots.harddisk.service_name
  operating_system = "byolinux_64"
  customizations {
   	hostname            = "mon-tux2"
   	http_headers = {
   		Authorization = "Basic bG9naW46cGFzc3dvcmQ="
   	}
   	image_check_sum  = "6a76c38a97f7a0909e6f53801dbe27b04c348a990592e075e910d65c3d45f1aa36c15a0530912b5c287a8151fc2891ec481d15eeab97ae7c918234abedb9c0a9"
    image_check_sum_type = "sha512"
    image_url          = "https://github.com/ashmonger/akution_test/releases/download/0.5-compress/deb11k6.qcow2"
  }
  properties = {
  	essential = "false"
	role      = "webservers"
  }
}
`

const testAccDedicatedServerReinstallConfig_RebootOnDestroy = `
data ovh_dedicated_server_boots "harddisk" {
  service_name = "%s"
  boot_type    = "harddisk"
}

data ovh_dedicated_server_boots "rescue" {
  service_name = data.ovh_dedicated_server_boots.harddisk.service_name
  boot_type    = "rescue"
}

resource ovh_dedicated_server_update "server" {
  service_name = data.ovh_dedicated_server_boots.harddisk.service_name
  boot_id      = data.ovh_dedicated_server_boots.harddisk.result[0]
  monitoring   = true
  state        = "ok"
}

resource ovh_dedicated_server_reinstall_task "server_reinstall" {
  service_name      = data.ovh_dedicated_server_boots.harddisk.service_name
  operating_system     = "debian12_64"
  bootid_on_destroy = data.ovh_dedicated_server_boots.rescue.result[0]
}
`

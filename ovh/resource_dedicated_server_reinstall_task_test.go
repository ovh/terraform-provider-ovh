package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDedicatedServerInstall_basic(t *testing.T) {
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
				Config: testAccDedicatedServerInstallConfig("basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "monitoring", "true"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_reinstall_task.server_install", "function", "reinstallServer"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_reinstall_task.server_install", "status", "done"),
				),
			},
		},
	})
}

func TestAccDedicatedServerInstall_rebootondestroy(t *testing.T) {
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
				Config: testAccDedicatedServerInstallConfig("rebootondestroy"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "monitoring", "true"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_reinstall_task.server_install", "function", "reinstallServer"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_reinstall_task.server_install", "status", "done"),
				),
			},
		},
	})
}

func testAccDedicatedServerInstallConfig(config string) string {
	dedicated_server := os.Getenv("OVH_DEDICATED_SERVER")
	sshKey := os.Getenv("OVH_SSH_KEY")
	if sshKey == "" {
		sshKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIrODOo0SvY5f0TlQNvGHIRKzr4bHPa+D5bYF18RiOgP email@example.com"
	}

	if config == "rebootondestroy" {
		return fmt.Sprintf(
			testAccDedicatedServerInstallConfig_RebootOnDestroy,
			dedicated_server,
		)
	}

	return fmt.Sprintf(
		testAccDedicatedServerInstallConfig_Basic,
		dedicated_server,
	)

}

const testAccDedicatedServerInstallConfig_Basic = `
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

resource ovh_dedicated_server_reinstall_task "server_install" {
  service_name = data.ovh_dedicated_server_boots.harddisk.service_name
  operating_system = "debian12_64"
}
`

const testAccDedicatedServerInstallConfig_RebootOnDestroy = `
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

resource ovh_dedicated_server_reinstall_task "server_install" {
  service_name      = data.ovh_dedicated_server_boots.harddisk.service_name
  operating_system     = "debian12_64"
  bootid_on_destroy = data.ovh_dedicated_server_boots.rescue.result[0]
}
`

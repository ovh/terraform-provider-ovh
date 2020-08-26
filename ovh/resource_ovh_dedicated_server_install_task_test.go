package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDedicatedServerInstall_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCredentials(t)
			testAccPreCheckDedicatedServer(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerInstallConfig("basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "monitoring", "true"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_install_task.server_install", "function", "reinstallServer"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_install_task.server_install", "status", "done"),
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
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerInstallConfig("rebootondestroy"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "monitoring", "true"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_install_task.server_install", "function", "reinstallServer"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_install_task.server_install", "status", "done"),
				),
			},
		},
	})
}

func testAccDedicatedServerInstallConfig(config string) string {
	dedicated_server := os.Getenv("OVH_DEDICATED_SERVER")
	testName := acctest.RandomWithPrefix(test_prefix)
	sshKey := os.Getenv("OVH_SSH_KEY")
	if sshKey == "" {
		sshKey = "ssh-ed25519 AAAAC3NzaC1yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"
	}

	if config == "rebootondestroy" {
		return fmt.Sprintf(
			testAccDedicatedServerInstallConfig_RebootOnDestroy,
			dedicated_server,
			testName,
			sshKey,
			testName,
		)
	}

	return fmt.Sprintf(
		testAccDedicatedServerInstallConfig_Basic,
		dedicated_server,
		testName,
		sshKey,
		testName,
	)

}

const testAccDedicatedServerInstallConfig_Basic = `
data ovh_dedicated_server_boots "harddisk" {
  service_name = "%s"
  boot_type    = "harddisk"
}

resource "ovh_me_ssh_key" "key" {
	key_name = "%s"
	key      = "%s"
}

resource ovh_dedicated_server_update "server" {
  service_name = data.ovh_dedicated_server_boots.harddisk.service_name
  boot_id      = data.ovh_dedicated_server_boots.harddisk.result[0]
  monitoring   = true
  state        = "ok"
}

resource "ovh_me_installation_template" "debian" {
  base_template_name = "debian10_64"
  template_name      = "%s"
  default_language   = "en"

  customization {
     change_log                      = "v1"
     custom_hostname                 = "mytest"
     ssh_key_name                    = ovh_me_ssh_key.key.key_name
  }
}

resource ovh_dedicated_server_install_task "server_install" {
  service_name = data.ovh_dedicated_server_boots.harddisk.service_name
  template_name = ovh_me_installation_template.debian.template_name
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

resource "ovh_me_ssh_key" "key" {
	key_name = "%s"
	key      = "%s"
}

resource ovh_dedicated_server_update "server" {
  service_name = data.ovh_dedicated_server_boots.harddisk.service_name
  boot_id      = data.ovh_dedicated_server_boots.harddisk.result[0]
  monitoring   = true
  state        = "ok"
}

resource "ovh_me_installation_template" "debian" {
  base_template_name = "debian10_64"
  template_name      = "%s"
  default_language   = "en"

  customization {
     change_log                      = "v1"
     custom_hostname                 = "mytest"
     ssh_key_name                    = ovh_me_ssh_key.key.key_name
  }
}

resource ovh_dedicated_server_install_task "server_install" {
  service_name      = data.ovh_dedicated_server_boots.harddisk.service_name
  template_name     = ovh_me_installation_template.debian.template_name
  bootid_on_destroy = data.ovh_dedicated_server_boots.rescue.result[0]
}
`

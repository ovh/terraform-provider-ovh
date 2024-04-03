package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
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

func TestAccDedicatedServerInstall_usermetadata(t *testing.T) {
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
				Config: testAccDedicatedServerInstallConfig("usermetadata"),
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

	if config == "usermetadata" {
		return fmt.Sprintf(
			testAccDedicatedServerInstallConfig_Usermetadata,
			dedicated_server,
			testName,
			sshKey,
			testName,
			sshKey,
			sshKey,
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
     custom_hostname                 = "mytest"
     ssh_key_name                    = ovh_me_ssh_key.key.key_name
  }
}

resource "time_sleep" "wait_for_ssh_key_sync" {
  create_duration = "120s"
  depends_on = [ovh_me_installation_template.debian]
}

resource ovh_dedicated_server_install_task "server_install" {
  service_name = data.ovh_dedicated_server_boots.harddisk.service_name
  template_name = ovh_me_installation_template.debian.template_name

  depends_on = [time_sleep.wait_for_ssh_key_sync]
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
  base_template_name = "debian12_64"
  template_name      = "%s"
  default_language   = "en"

  customization {
     custom_hostname                 = "mytest"
     ssh_key_name                    = ovh_me_ssh_key.key.key_name
  }
}

resource "time_sleep" "wait_for_ssh_key_sync" {
	create_duration = "120s"
	depends_on = [ovh_me_installation_template.debian]
  }

resource ovh_dedicated_server_install_task "server_install" {
  service_name      = data.ovh_dedicated_server_boots.harddisk.service_name
  template_name     = ovh_me_installation_template.debian.template_name
  bootid_on_destroy = data.ovh_dedicated_server_boots.rescue.result[0]

  depends_on = [time_sleep.wait_for_ssh_key_sync]
}
`
const testAccDedicatedServerInstallConfig_Usermetadata = `
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
  monitoring   = true
  state        = "ok"
}

resource "ovh_me_installation_template" "byolinux" {
  base_template_name = "byolinux_64"
  template_name      = "%s"
  default_language   = "en"
}

resource ovh_dedicated_server_install_task "server_install" {
  service_name      = data.ovh_dedicated_server_boots.harddisk.service_name
  template_name     = ovh_me_installation_template.byolinux.template_name
  user_metadata {
	key   = "imageURL"
	value = "https://github.com/ashmonger/akution_test/releases/download/0.6-fixCache/deb11k6.qcow2"
  }
  user_metadata {
	key   = "imageType"
	value = "qcow2"
  }
  user_metadata {
	key   = "httpHeaders0Key"
	value = "Authorization"
  }
  user_metadata {
	key   = "httpHeaders0Value"
	value = "Basic bG9naW46cGFzc3dvcmQ="
  }
  user_metadata {
	key   = "imageCheckSum"
	value = "047122c9ff4d2a69512212104b06c678f5a9cdb22b75467353613ff87ccd03b57b38967e56d810e61366f9d22d6bd39ac0addf4e00a4c6445112a2416af8f225"
  }
  user_metadata {
	key   = "imageCheckSumType"
	value = "sha512"
  }
  user_metadata {
	key   = "configDriveUserData"
	value = "#cloud-config\nssh_authorized_keys:\n  - %s\n\nusers:\n  - name: aautret\n    sudo: ALL=(ALL) NOPASSWD:ALL\n    groups: users, sudo\n    shell: /bin/bash\n    lock_passwd: false\n    ssh_authorized_keys:\n      - %s\ndisable_root: false\npackages:\n  - vim\n  - tree\nfinal_message: The system is finally up, after $UPTIME seconds\n"
  }
}
`

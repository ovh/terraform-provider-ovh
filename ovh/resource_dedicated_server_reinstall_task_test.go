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
						"ovh_dedicated_server_update.server", "monitoring", "false"),
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
						"ovh_dedicated_server_update.server", "monitoring", "false"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_reinstall_task.server_reinstall", "function", "reinstallServer"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_reinstall_task.server_reinstall", "status", "done"),
				),
			},
		},
	})
}

func TestAccDedicatedServerReinstall_customizations(t *testing.T) {
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
				Config: testAccDedicatedServerReinstallConfig("customizations"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "monitoring", "false"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_reinstall_task.server_reinstall", "function", "reinstallServer"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_reinstall_task.server_reinstall", "status", "done"),
				),
			},
		},
	})
}

func TestAccDedicatedServerReinstall_byolinux(t *testing.T) {
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
				Config: testAccDedicatedServerReinstallConfig("byolinux"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "monitoring", "false"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_reinstall_task.server_reinstall", "function", "reinstallServer"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_reinstall_task.server_reinstall", "status", "done"),
				),
			},
		},
	})
}

func TestAccDedicatedServerReinstall_storage(t *testing.T) {
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
				Config: testAccDedicatedServerReinstallConfig("storage"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "monitoring", "false"),
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

	if config == "customizations" {
		return fmt.Sprintf(
			testAccDedicatedServerReinstallConfig_Customizations,
			dedicated_server,
		)
	}

	if config == "byolinux" {
		return fmt.Sprintf(
			testAccDedicatedServerReinstallConfig_Byolinux,
			dedicated_server,
		)
	}

	if config == "storage" {
		return fmt.Sprintf(
			testAccDedicatedServerReinstallConfig_Storage,
			dedicated_server,
		)
	}

	return fmt.Sprintf(
		testAccDedicatedServerReinstallConfig_Basic,
		dedicated_server,
	)

}

const testAccDedicatedServerReinstallConfig_Basic = `
data "ovh_dedicated_server_boots" "harddisk" {
  service_name = "%s"
  boot_type    = "harddisk"
}

resource "ovh_dedicated_server_update" "server" {
  service_name        = data.ovh_dedicated_server_boots.harddisk.service_name
  boot_id             = data.ovh_dedicated_server_boots.harddisk.result[0]
  monitoring          = false
  state               = "ok"
  efi_bootloader_path = "\\efi\\debian\\grubx64.efi"
}

resource "ovh_dedicated_server_reinstall_task" "server_reinstall" {
  service_name     = data.ovh_dedicated_server_boots.harddisk.service_name
  os= "debian12_64"
}
`

const testAccDedicatedServerReinstallConfig_RebootOnDestroy = `
data "ovh_dedicated_server_boots" "harddisk" {
  service_name = "%s"
  boot_type    = "harddisk"
}

data "ovh_dedicated_server_boots" "rescue" {
  service_name = data.ovh_dedicated_server_boots.harddisk.service_name
  boot_type    = "rescue"
}

resource "ovh_dedicated_server_update" "server" {
  service_name        = data.ovh_dedicated_server_boots.harddisk.service_name
  boot_id             = data.ovh_dedicated_server_boots.harddisk.result[0]
  monitoring          = false
  state               = "ok"
  efi_bootloader_path = "\\efi\\debian\\grubx64.efi"
}

resource "ovh_dedicated_server_reinstall_task" "server_reinstall" {
  service_name      = data.ovh_dedicated_server_boots.harddisk.service_name
  os     = "debian12_64"
  bootid_on_destroy = data.ovh_dedicated_server_boots.rescue.result[0]
}
`

const testAccDedicatedServerReinstallConfig_Customizations = `
data "ovh_dedicated_server_boots" "harddisk" {
  service_name = "%s"
  boot_type    = "harddisk"
}

resource "ovh_dedicated_server_update" "server" {
  service_name = data.ovh_dedicated_server_boots.harddisk.service_name
  boot_id      = data.ovh_dedicated_server_boots.harddisk.result[0]
  monitoring   = false
  state        = "ok"
}

resource "ovh_dedicated_server_reinstall_task" "server_reinstall" {
  service_name     = data.ovh_dedicated_server_boots.harddisk.service_name
  os = "debian12_64"
  customizations {
    hostname               = "mon-tux"
    post_installation_script = "IyEvYmluL2Jhc2gKZWNobyAiY291Y291IHBvc3RJbnN0YWxsYXRpb25TY3JpcHQiID4gL29wdC9jb3Vjb3UKY2F0IC9ldGMvbWFjaGluZS1pZCAgPj4gL29wdC9jb3Vjb3UKZGF0ZSAiKyVZLSVtLSVkICVIOiVNOiVTIiAtLXV0YyA+PiAvb3B0L2NvdWNvdQo="
    ssh_key                 = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQC9xPpdqP3sx2H+gcBm65tJEaUbuifQ1uGkgrWtNY0PRKNNPdy+3yoVOtxk6Vjo4YZ0EU/JhmQfnrK7X7Q5vhqYxmozi0LiTRt0BxgqHJ+4hWTWMIOgr+C2jLx7ZsCReRk+fy5AHr6h0PHQEuXVLXeUy/TDyuY2JPtUZ5jcqvLYgQ== my-nuclear-power-plant"
  }
}
`
const testAccDedicatedServerReinstallConfig_Byolinux = `
data "ovh_dedicated_server_boots" "harddisk" {
  service_name = "%s"
  boot_type    = "harddisk"
}

resource "ovh_dedicated_server_update" "server" {
  service_name = data.ovh_dedicated_server_boots.harddisk.service_name
  boot_id      = data.ovh_dedicated_server_boots.harddisk.result[0]
  monitoring   = false
  state        = "ok"
}

resource "ovh_dedicated_server_reinstall_task" "server_reinstall" {
  service_name     = data.ovh_dedicated_server_boots.harddisk.service_name
  os = "byolinux_64"
  customizations {
    hostname = "mon-tux"
    image_url = "https://github.com/ashmonger/akution_test/releases/latest/download/deb11k6.qcow2"
	efi_bootloader_path = "\\efi\\debian\\grubx64.efi"
	ssh_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIrODOo0SvY5f0TlQNvGHIRKzr4bHPa+D5bYF18RiOgP email@example.com"
	config_drive_user_data = "c3NoX2F1dGhvcml6ZWRfa2V5czoKICAtIHNzaC1yc2EgQUFBQUI4ZGpZaXc9PSBteXNlbGZAbXlkb21haW4ubmV0Cgp1c2VyczoKICAtIG5hbWU6IHBhdGllbnQwCiAgICBzdWRvOiBBTEw9KEFMTCkgTk9QQVNTV0Q6QUxMCiAgICBncm91cHM6IHVzZXJzLCBzdWRvCiAgICBzaGVsbDogL2Jpbi9iYXNoCiAgICBsb2NrX3Bhc3N3ZDogZmFsc2UKICAgIHNzaF9hdXRob3JpemVkX2tleXM6CiAgICAgIC0gc3NoLXJzYSBBQUFBQjhkallpdz09IG15c2VsZkBteWRvbWFpbi5uZXQKZGlzYWJsZV9yb290OiBmYWxzZQpwYWNrYWdlczoKICAtIHZpbQogIC0gdHJlZQpmaW5hbF9tZXNzYWdlOiBUaGUgc3lzdGVtIGlzIGZpbmFsbHkgdXAsIGFmdGVyICRVUFRJTUUgc2Vjb25kcw=="
	config_drive_metadata = {
		foo = "bar"
		hello = "world"
    }
	http_headers = {
		Authorization = "Basic bG9naW46cGFzc3dvcmQ="
	}
  }
}
`

const testAccDedicatedServerReinstallConfig_Storage = `
data "ovh_dedicated_server_boots" "harddisk" {
  service_name = "%s"
  boot_type    = "harddisk"
}

resource "ovh_dedicated_server_update" "server" {
  service_name = data.ovh_dedicated_server_boots.harddisk.service_name
  boot_id      = data.ovh_dedicated_server_boots.harddisk.result[0]
  monitoring   = false
  state        = "ok"
}

resource "ovh_dedicated_server_reinstall_task" "server_reinstall" {
  service_name     = data.ovh_dedicated_server_boots.harddisk.service_name
  os = "debian12_64"
  customizations {
    hostname = "mon-tux"
  }
  storage {
    partitioning {
      disks = 2
      layout {
        file_system = "ext4"
        mount_point = "/boot"
        raid_level  = 1
        size       = 1024
      }
      layout {
        file_system = "ext4"
        mount_point = "/"
        raid_level  = 1
        size       = 20480
        extras {
          lv {
            name = "root"
          }
        }
      }
      layout {
        file_system = "swap"
        mount_point = "swap"
        size       = 2048
      }
      layout {
        file_system = "zfs"
        mount_point = "/data"
        raid_level  = 5
        size       = 0
        extras {
          zp {
            name = "poule"
          }
        }
      }
    }
  }
}
`

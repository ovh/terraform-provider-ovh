data "ovh_dedicated_server" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

data "ovh_dedicated_installation_template" "template" {
  template_name = "debian12_64"
}

resource "ovh_dedicated_server_reinstall_task" "server_install" {
  service_name     = data.ovh_dedicated_server.server.service_name
  os = data.ovh_dedicated_installation_template.template.template_name
  customizations {
    hostname = "mon-tux"
  }
  storage {
    partitioning {
      disk_group_id = 2
      hardware_raid {
        raid_level = 5
      }
      layout {
        file_system = "ext4"
        mount_point = "/"
        raid_level  = 1
        size        = 20480
      }
    }
  }
}

data "ovh_dedicated_server" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

data "ovh_dedicated_installation_template" "template" {
  template_name = "debian12_64"
}

resource "ovh_dedicated_server_reinstall_task" "server_install" {
  service_name = data.ovh_dedicated_server.server.service_name
  os           = data.ovh_dedicated_installation_template.template.template_name
  customizations {
    hostname = "mon-tux"
  }
  storage {
    disk_group_id = 1
    partitioning {
      scheme_name = "default"
    }
  }
  storage {
    disk_group_id = 2
    erase         = false
  }
}

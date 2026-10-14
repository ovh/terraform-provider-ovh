data "ovh_dedicated_server" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

data "ovh_dedicated_installation_template" "template" {
  template_name = "ubuntu2404-withapp_64"
}

locals {
  app_to_install_input = one([
    for input in data.ovh_dedicated_installation_template.template.inputs : input
    if input.name == "appToInstall"
  ])
}

# List of applications available for this OS
output "available_apps" {
  value = local.app_to_install_input.enum
}

resource "ovh_dedicated_server_reinstall_task" "server_reinstall" {
  service_name = data.ovh_dedicated_server.server.service_name
  os           = data.ovh_dedicated_installation_template.template.template_name
  customizations {
    hostname       = "mon-tux"
    app_to_install = "wordpress"
  }
}

data "ovh_dedicated_server" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

data "ovh_dedicated_installation_template" "template" {
  template_name = "debian12_64"
}

resource "ovh_dedicated_server_reinstall_task" "server_reinstall" {
  service_name = data.ovh_dedicated_server.server.service_name
  os           = data.ovh_dedicated_installation_template.template.template_name
  customizations {
    hostname                 = "mon-tux"
    post_installation_script = "IyEvYmluL2Jhc2gKZWNobyAiY291Y291IHBvc3RJbnN0YWxsYXRpb25TY3JpcHQiID4gL29wdC9jb3Vjb3UKY2F0IC9ldGMvbWFjaGluZS1pZCAgPj4gL29wdC9jb3Vjb3UKZGF0ZSAiKyVZLSVtLSVkICVIOiVNOiVTIiAtLXV0YyA+PiAvb3B0L2NvdWNvdQo="
    ssh_key                  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQC9xPpdqP3sx2H+gcBm65tJEaUbuifQ1uGkgrWtNY0PRKNNPdy+3yoVOtxk6Vjo4YZ0EU/JhmQfnrK7X7Q5vhqYxmozi0LiTRt0BxgqHJ+4hWTWMIOgr+C2jLx7ZsCReRk+fy5AHr6h0PHQEuXVLXeUy/TDyuY2JPtUZ5jcqvLYgQ== my-nuclear-power-plant"
  }
}

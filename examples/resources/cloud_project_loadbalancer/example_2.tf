resource "ovh_cloud_project_network_private" "priv" {
  service_name  = "<public cloud project ID>"
  vlan_id       = "10"
  name          = "my_priv"
  regions       = ["GRA9"]
}

resource "ovh_cloud_project_network_private_subnet" "privsub" {
  service_name  = ovh_cloud_project_network_private.priv.service_name
  network_id    = ovh_cloud_project_network_private.priv.id
  region        = "GRA9"
  start         = "10.0.0.2"
  end           = "10.0.255.254"
  network       = "10.0.0.0/16"
  dhcp          = true
}

resource "ovh_cloud_project_loadbalancer" "lb" {
  service_name = ovh_cloud_project_network_private_subnet.privsub.service_name
  region_name = ovh_cloud_project_network_private_subnet.privsub.region
  flavor_id = "<loadbalancer flavor ID>"
  network = {
    private = {
      network = {
        id = element([for region in ovh_cloud_project_network_private.priv.regions_attributes: region if "${region.region}" == "GRA9"], 0).openstackid
        subnet_id = ovh_cloud_project_network_private_subnet.privsub.id
      }
    }
  }
  description = "My new LB"
  listeners = [
    {
      port = "34568"
      protocol = "tcp"
    },
    {
      port = "34569"
      protocol = "udp"
    }
  ]
}

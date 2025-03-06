resource "ovh_cloud_project_loadbalancer" "lb" {
  service_name = "<public cloud project ID>"
  region_name = "GRA9"
  flavor_id = "<loadbalancer flavor ID>"
  network = {
    private = {
      network = {
        id = element([for region in ovh_cloud_project_network_private.mypriv.regions_attributes: region if "${region.region}" == "GRA9"], 0).openstackid
        subnet_id = ovh_cloud_project_network_private_subnet.myprivsub.id
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

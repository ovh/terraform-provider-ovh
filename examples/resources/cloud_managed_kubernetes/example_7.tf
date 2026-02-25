
resource "ovh_cloud_project_network_private" "network" {
  service_name = var.service_name # Public Cloud service name
  vlan_id     = 42
  name       = "terraform_testacc_private_net"
  regions    = ["GRA11"]
}

resource "ovh_cloud_project_network_private_subnet" "subnet" {
  service_name = var.service_name
  network_id   = ovh_cloud_project_network_private.network.id

  # whatever region, for test purpose
  region     = "GRA11"
  start      = "192.168.168.100"
  end        = "192.168.168.200"
  network    = "192.168.168.0/24"
  dhcp       = true
  no_gateway = false
}

resource "ovh_cloud_project_gateway" "gateway" {
  service_name = var.service_name
  name       = "gateway"
  model      = "s"
  region     = "GRA11"
  network_id = tolist(ovh_cloud_project_network_private.network.regions_attributes[*].openstackid)[0]
  subnet_id  = ovh_cloud_project_network_private_subnet.subnet.id
}

resource "ovh_cloud_managed_kubernetes" "my_cluster" {
  service_name  = var.service_name
  name          = "test-kube-attach"
  region        = "GRA11"

  private_network_id = tolist(ovh_cloud_project_network_private.network.regions_attributes[*].openstackid)[0]
  nodes_subnet_id = ovh_cloud_project_network_private_subnet.subnet.id
  private_network_configuration {
      default_vrack_gateway              = ""
      private_network_routing_as_default = false
  }
}

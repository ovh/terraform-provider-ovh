resource "ovh_cloud_network_private_vrack" "network" {
  service_name = <Public cloud project id>
  name         = "my-network"
  region       = "GRA1"
}

resource "ovh_cloud_network_private_vrack_subnet" "subnet" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id   = ovh_cloud_network_private_vrack.network.id
  name         = "my-subnet"
  cidr         = "10.0.0.0/24"
  region       = "GRA1"
}

resource "ovh_cloud_loadbalancer" "lb" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  name         = "my-loadbalancer"
  region       = "GRA1"
  flavor_name  = "SMALL"
  description  = "My load balancer"

  network = {
    id        = ovh_cloud_network_private_vrack.network.id
    subnet_id = ovh_cloud_network_private_vrack_subnet.subnet.id
  }
}

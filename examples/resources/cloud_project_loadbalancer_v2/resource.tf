resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "xxxxxxxxxx"
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
  service_name   = ovh_cloud_network_private_vrack.network.service_name
  name           = "my-loadbalancer"
  region         = "GRA1"
  vip_network_id = ovh_cloud_network_private_vrack.network.id
  vip_subnet_id  = ovh_cloud_network_private_vrack_subnet.subnet.id
  flavor_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  description    = "My load balancer"
}

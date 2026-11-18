resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name         = "databases"
  region       = "EU-WEST-PAR"
}

resource "ovh_cloud_network_private_vrack_subnet" "subnet" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id   = ovh_cloud_network_private_vrack.network.id
  name         = "databases"
  cidr         = "10.0.0.0/24"
  region       = ovh_cloud_network_private_vrack.network.region
}

resource "ovh_cloud_managed_database" "db" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  engine       = "postgresql"
  description  = "my-private-postgresql"
  region       = "eu-west-par"
  version      = "17"
  plan         = "business"
  flavor       = "b3-8"
  nodes_count  = 2

  network = {
    network_id = ovh_cloud_network_private_vrack.network.id
    subnet_id  = ovh_cloud_network_private_vrack_subnet.subnet.id
  }
}

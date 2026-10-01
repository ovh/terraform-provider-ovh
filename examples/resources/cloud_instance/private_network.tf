resource "ovh_cloud_network_private_vrack" "private" {
  service_name = "<Public cloud project id>"
  name         = "my-private-network"
  region       = "GRA11"
}

resource "ovh_cloud_network_private_vrack_subnet" "private" {
  service_name = ovh_cloud_network_private_vrack.private.service_name
  network_id   = ovh_cloud_network_private_vrack.private.id
  name         = "my-subnet"
  cidr         = "10.0.0.0/24"
  region       = "GRA11"
}

# No public interface at all: the instance is only reachable from the vRack.
# Add an ovh_cloud_gateway on the same subnet to give it egress.
resource "ovh_cloud_instance" "private_only" {
  service_name = ovh_cloud_network_private_vrack.private.service_name
  region       = "GRA11"
  name         = "my-private-instance"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"

  networks = [
    {
      network_id = ovh_cloud_network_private_vrack.private.id
      subnet_id  = ovh_cloud_network_private_vrack_subnet.private.id
    },
  ]
}

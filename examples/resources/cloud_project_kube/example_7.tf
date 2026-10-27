resource "ovh_cloud_network_private_vrack" "network" {
  service_name = var.service_name # Public Cloud service name
  name         = "terraform_testacc_private_net"
  region       = "GRA11"
  vlan_id      = 42
}

resource "ovh_cloud_network_private_vrack_subnet" "subnet" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id   = ovh_cloud_network_private_vrack.network.id
  name         = "terraform_testacc_private_subnet"
  region       = "GRA11"
  cidr         = "192.168.168.0/24"
  gateway_ip   = "192.168.168.1"
  dhcp_enabled = true
  allocation_pools = [{
    start = "192.168.168.100"
    end   = "192.168.168.200"
  }]
}

resource "ovh_cloud_gateway" "gateway" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  name         = "gateway"
  region       = "GRA11"
  external_gateway = {
    enabled = true
    model   = "S"
  }
  subnet_ids = [ovh_cloud_network_private_vrack_subnet.subnet.id]
}

resource "ovh_cloud_project_kube" "my_cluster" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  name         = "test-kube-attach"
  region       = "GRA11"

  private_network_id = ovh_cloud_network_private_vrack.network.id
  nodes_subnet_id    = ovh_cloud_network_private_vrack_subnet.subnet.id
  private_network_configuration {
    default_vrack_gateway              = ""
    private_network_routing_as_default = true
  }

  depends_on = [ovh_cloud_gateway.gateway]
}

# All nodes reach the internet through the gateway: use this IP to allow them on external services
output "egress_ip" {
  value = ovh_cloud_gateway.gateway.current_state.external_ip
}

resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" # Public Cloud service name
  name         = "terraform_mks_multiaz_private_net"
  region       = "EU-WEST-PAR"
  vlan_id      = 84
}

resource "ovh_cloud_network_private_vrack_subnet" "subnet" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id   = ovh_cloud_network_private_vrack.network.id
  name         = "terraform_mks_multiaz_private_subnet"
  region       = "EU-WEST-PAR"
  cidr         = "192.168.142.0/24"
  gateway_ip   = "192.168.142.1"
  dhcp_enabled = true
  allocation_pools = [{
    start = "192.168.142.100"
    end   = "192.168.142.200"
  }]
}

resource "ovh_cloud_gateway" "gateway" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  name         = "gateway"
  region       = "EU-WEST-PAR"
  external_gateway = {
    enabled = true
    model   = "S"
  }
  subnet_ids = [ovh_cloud_network_private_vrack_subnet.subnet.id]
}

resource "ovh_cloud_project_kube" "my_multizone_cluster" {
  service_name  = ovh_cloud_network_private_vrack.network.service_name
  name          = "multi-zone-mks"
  region        = "EU-WEST-PAR"
  plan          = "standard"

  private_network_id = ovh_cloud_network_private_vrack.network.id
  nodes_subnet_id    = ovh_cloud_network_private_vrack_subnet.subnet.id

  depends_on    = [ ovh_cloud_gateway.gateway ] //Gateway is mandatory for multizones cluster
}

resource "ovh_cloud_project_kube_nodepool" "node_pool_multi_zones_a" {
  service_name       = ovh_cloud_network_private_vrack.network.service_name
  kube_id            = ovh_cloud_project_kube.my_multizone_cluster.id
  name               = "my-pool-zone-a" //Warning: "_" char is not allowed!
  flavor_name        = "b3-8"
  desired_nodes      = 3
  availability_zones = ["eu-west-par-a"] //Currently, only one zone is supported
}

resource "ovh_cloud_project_kube_nodepool" "node_pool_multi_zones_b" {
  service_name       = ovh_cloud_network_private_vrack.network.service_name
  kube_id            = ovh_cloud_project_kube.my_multizone_cluster.id
  name               = "my-pool-zone-b"
  flavor_name        = "b3-8"
  desired_nodes      = 3
  availability_zones = ["eu-west-par-b"]
}

resource "ovh_cloud_project_kube_nodepool" "node_pool_multi_zones_c" {
  service_name       = ovh_cloud_network_private_vrack.network.service_name
  kube_id            = ovh_cloud_project_kube.my_multizone_cluster.id
  name               = "my-pool-zone-c"
  flavor_name        = "b3-8"
  desired_nodes      = 3
  availability_zones = ["eu-west-par-c"]
}

output "kubeconfig_file_eu_west_par" {
  value     = ovh_cloud_project_kube.my_multizone_cluster.kubeconfig
  sensitive = true
}

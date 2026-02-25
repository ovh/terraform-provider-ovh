resource "ovh_cloud_managed_kubernetes" "my_multizone_cluster" {
  service_name       = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name               = "terraform-multi-zone-cluster"
  region             = "EU-WEST-PAR"
  plan               = "standard"

  private_network_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  nodes_subnet_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" //Subnet must has a OVHcloud Gateway/OpenStack router
}

resource "ovh_cloud_managed_kubernetes_nodepool" "node_pool_multi_zones_a" {
  service_name       = var.service_name
  kube_id            = ovh_cloud_managed_kubernetes.my_multizone_cluster.id
  name               = "my-pool-zone-a" //Warning: "_" char is not allowed!
  flavor_name        = "b3-8"
  desired_nodes      = 3
  availability_zones = ["eu-west-par-a"] //Currently, only one zone is supported
}

resource "ovh_cloud_managed_kubernetes_nodepool" "node_pool_multi_zones_b" {
  service_name       = var.service_name
  kube_id            = ovh_cloud_managed_kubernetes.my_multizone_cluster.id
  name               = "my-pool-zone-b"
  flavor_name        = "b3-8"
  desired_nodes      = 3
  availability_zones = ["eu-west-par-b"]
}

resource "ovh_cloud_managed_kubernetes_nodepool" "node_pool_multi_zones_c" {
  service_name       = var.service_name
  kube_id            = ovh_cloud_managed_kubernetes.my_multizone_cluster.id
  name               = "my-pool-zone-c"
  flavor_name        = "b3-8"
  desired_nodes      = 3
  availability_zones = ["eu-west-par-c"]
}

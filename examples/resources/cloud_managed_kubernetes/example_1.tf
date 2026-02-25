resource "ovh_cloud_managed_kubernetes" "my_cluster" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name         = "my_kube_cluster"
  region       = "GRA11"
}

resource "ovh_cloud_managed_kubernetes_nodepool" "node_pool_1" {
  service_name  = ovh_cloud_managed_kubernetes.my_cluster.service_name
  kube_id       = ovh_cloud_managed_kubernetes.my_cluster.id
  name          = "my-pool-1"
  flavor_name   = "b3-8"
  desired_nodes = 3
}

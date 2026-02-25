resource "ovh_cloud_managed_kubernetes_nodepool" "node_pool" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  kube_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name          = "my-pool-1" //Warning: "_" char is not allowed!
  flavor_name   = "b3-8"
  desired_nodes = 3
}

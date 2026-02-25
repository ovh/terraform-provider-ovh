data "ovh_cloud_managed_kubernetes_nodepool" "nodepool" {
  service_name  = "XXXXXX"
  kube_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxx"
  name          = "xxxxxx"
}

output "max_nodes" {
  value = data.ovh_cloud_managed_kubernetes_nodepool.nodepool.max_nodes
}

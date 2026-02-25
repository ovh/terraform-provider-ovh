data "ovh_cloud_managed_kubernetes_nodepool_nodes" "nodes" {
  service_name  = "XXXXXX"
  kube_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxx"
  name          = "XXXXXX"
}

output "nodes" {
  value = data.ovh_cloud_managed_kubernetes_nodepool_nodes.nodes
}

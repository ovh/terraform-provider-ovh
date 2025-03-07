data "ovh_cloud_project_kube_nodepool" "nodepool" {
  service_name  = "XXXXXX"
  kube_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxx"
  name          = "xxxxxx"
}

output "max_nodes" {
  value = data.ovh_cloud_project_kube_nodepool.nodepool.max_nodes
}

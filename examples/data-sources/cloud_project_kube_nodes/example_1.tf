data "ovh_cloud_project_kube_nodes" "nodes" {
  service_name  = "XXXXXX"
  kube_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxx"
}

output "nodes" {
  value = data.ovh_cloud_project_kube_nodes.nodes
}

data "ovh_cloud_project_kube" "my_kube_cluster" {
  service_name = "XXXXXX"
  kube_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}

output "version" {
  value = data.ovh_cloud_project_kube.my_kube_cluster.version
}

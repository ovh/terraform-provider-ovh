data "ovh_cloud_managed_kubernetes_versions" "versions" {
  service_name = "<Public cloud project id>"
}

output "default_version" {
  value = data.ovh_cloud_managed_kubernetes_versions.versions.default_version
}

data "ovh_cloud_project_network_privates" "private" {
  service_name = "XXXXXX"
}

output "private" {
  value = data.ovh_cloud_project_network_privates.private
}

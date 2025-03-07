data "ovh_cloud_project_network_private" "private" {
  service_name = "XXXXXX"
  network_id           = "XXX"
}
output "private" {
  value = data.ovh_cloud_project_network_private.private
}

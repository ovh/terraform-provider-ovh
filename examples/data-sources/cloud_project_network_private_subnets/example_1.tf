data "ovh_cloud_project_network_private_subnets" "private" {
  service_name = "XXXXXX"
  network_id   = "XXXXXX"
}
output "private" {
  value = data.ovh_cloud_project_network_private_subnets.private
}

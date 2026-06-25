data "ovh_cloud_key_manager_container" "container" {
  service_name = "Public cloud project ID"
  container_id = "00000000-0000-0000-0000-000000000000"
}

output "container_name" {
  value = data.ovh_cloud_key_manager_container.container.name
}

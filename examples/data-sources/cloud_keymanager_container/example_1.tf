data "ovh_cloud_keymanager_container" "container" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  container_id = "00000000-0000-0000-0000-000000000000"
}

output "container_name" {
  value = data.ovh_cloud_keymanager_container.container.name
}

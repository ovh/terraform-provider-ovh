data "ovh_cloud_keymanager_containers" "all" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
}

output "container_ids" {
  value = [for c in data.ovh_cloud_keymanager_containers.all.containers : c.id]
}

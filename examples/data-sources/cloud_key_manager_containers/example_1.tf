data "ovh_cloud_key_manager_containers" "all" {
  service_name = "Public cloud project ID"
}

output "container_ids" {
  value = [for c in data.ovh_cloud_key_manager_containers.all.containers : c.id]
}

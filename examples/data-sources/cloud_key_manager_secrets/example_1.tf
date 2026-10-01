data "ovh_cloud_key_manager_secrets" "all" {
  service_name = "Public cloud project ID"
}

output "secret_ids" {
  value = [for s in data.ovh_cloud_key_manager_secrets.all.secrets : s.id]
}

data "ovh_cloud_key_manager_secret" "secret" {
  service_name = "Public cloud project ID"
  secret_id    = "00000000-0000-0000-0000-000000000000"
}

output "secret_name" {
  value = data.ovh_cloud_key_manager_secret.secret.name
}

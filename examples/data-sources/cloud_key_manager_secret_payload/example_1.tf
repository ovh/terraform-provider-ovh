data "ovh_cloud_key_manager_secret_payload" "payload" {
  service_name = "Public cloud project ID"
  secret_id    = "00000000-0000-0000-0000-000000000000"
}

output "secret_payload" {
  value     = data.ovh_cloud_key_manager_secret_payload.payload.payload
  sensitive = true
}

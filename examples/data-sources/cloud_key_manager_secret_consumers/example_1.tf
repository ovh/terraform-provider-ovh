data "ovh_cloud_key_manager_secret_consumers" "consumers" {
  service_name = "Public cloud project ID"
  secret_id    = "00000000-0000-0000-0000-000000000000"
}

output "consumer_services" {
  value = [for c in data.ovh_cloud_key_manager_secret_consumers.consumers.consumers : c.service]
}

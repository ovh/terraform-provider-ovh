data "ovh_cloud_keymanager_secret_consumers" "consumers" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  secret_id    = "00000000-0000-0000-0000-000000000000"
}

output "consumer_services" {
  value = [for c in data.ovh_cloud_keymanager_secret_consumers.consumers.consumers : c.service]
}

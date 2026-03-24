data "ovh_cloud_keymanager_secret" "secret" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  secret_id    = "00000000-0000-0000-0000-000000000000"
}

output "secret_name" {
  value = data.ovh_cloud_keymanager_secret.secret.name
}

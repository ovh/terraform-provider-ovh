data "ovh_cloud_keymanager_secrets" "all" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
}

output "secret_ids" {
  value = [for s in data.ovh_cloud_keymanager_secrets.all.secrets : s.id]
}

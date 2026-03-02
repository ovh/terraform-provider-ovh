data "ovh_cloud_managed_database_valkey_user" "valkey_user" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ"
}

output "valkey_user_commands" {
  value = data.ovh_cloud_managed_database_valkey_user.valkey_user.commands
}

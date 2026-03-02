data "ovh_cloud_managed_database_postgresql_user" "pg_user" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ"
}

output "pg_user_roles" {
  value = data.ovh_cloud_managed_database_postgresql_user.pg_user.roles
}

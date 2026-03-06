data "ovh_cloud_project_database_clickhouse_user" "ch_user" {
  service_name = "XXX"
  cluster_id   = "YYY"
  name         = "ZZZ"
}

output "ch_user_roles" {
  value = data.ovh_cloud_project_database_clickhouse_user.ch_user.roles
}

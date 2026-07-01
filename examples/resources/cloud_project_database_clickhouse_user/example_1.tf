data "ovh_cloud_project_database" "clickhouse" {
  service_name  = "XXXX"
  engine        = "clickhouse"
  id            = "ZZZZ"
}

resource "ovh_cloud_project_database_clickhouse_user" "user" {
  service_name  = data.ovh_cloud_project_database.clickhouse.service_name
  cluster_id    = data.ovh_cloud_project_database.clickhouse.id
  name          = "johndoe"
}

output "user_password" {
  value     = ovh_cloud_project_database_clickhouse_user.user.password
  sensitive = true
}

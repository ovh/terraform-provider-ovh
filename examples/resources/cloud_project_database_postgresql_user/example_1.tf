data "ovh_cloud_project_database" "postgresql" {
  service_name  = "XXXX"
  engine        = "postgresql"
  id            = "ZZZZ"
}

resource "ovh_cloud_project_database_postgresql_user" "user" {
  service_name  = data.ovh_cloud_project_database.postgresql.service_name
  cluster_id    = data.ovh_cloud_project_database.postgresql.id
  name          = "johndoe"
  roles         = ["replication"]
}

output "user_password" {
  value     = ovh_cloud_project_database_postgresql_user.user.password
  sensitive = true
}

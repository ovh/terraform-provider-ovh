resource "ovh_cloud_managed_database_database" "app" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  engine       = "postgresql"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  name         = "app"
}

resource "ovh_cloud_managed_database_postgresql_user" "app" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  name         = "app"
}

resource "ovh_cloud_managed_database_postgresql_connection_pool" "app" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  name         = "app"
  database_id  = ovh_cloud_managed_database_database.app.id
  user_id      = ovh_cloud_managed_database_postgresql_user.app.id
  mode         = "TRANSACTION"
  size         = 10
}

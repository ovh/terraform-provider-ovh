data "ovh_cloud_managed_database" "db" {
  service_name  = "XXXX"
  engine        = "YYYY"
  id            = "ZZZZ"
}

resource "ovh_cloud_managed_database_user" "user_datetime" {
  service_name    = data.ovh_cloud_managed_database.db.service_name
  engine          = data.ovh_cloud_managed_database.db.engine
  cluster_id      = data.ovh_cloud_managed_database.db.id
  name            = "alice"
  password_reset  = "2024-01-02T11:00:00Z"
}

resource "ovh_cloud_managed_database_user" "user_md5" {
  service_name    = data.ovh_cloud_managed_database.db.service_name
  engine          = data.ovh_cloud_managed_database.db.engine
  cluster_id      = data.ovh_cloud_managed_database.db.id
  name            = "bob"
  password_reset  = "md5(var.something)"
}

resource "ovh_cloud_managed_database_user" "user" {
  service_name    = data.ovh_cloud_managed_database.db.service_name
  engine          = data.ovh_cloud_managed_database.db.engine
  cluster_id      = data.ovh_cloud_managed_database.db.id
  name            = "johndoe"
  password_reset  = "reset1"
}

output "user_password" {
  value     = ovh_cloud_managed_database_user.user.password
  sensitive = true
}

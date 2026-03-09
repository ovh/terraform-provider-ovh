data "ovh_cloud_project_database" "clickhouse" {
  service_name  = "XXXX"
  engine        = "clickhouse"
  id            = "ZZZZ"
}

# Change password_reset with the datetime each time you want to reset the password to trigger an update
resource "ovh_cloud_project_database_clickhouse_user" "user_datetime" {
  service_name    = data.ovh_cloud_project_database.clickhouse.service_name
  cluster_id      = data.ovh_cloud_project_database.clickhouse.id
  name            = "alice"
  password_reset  = "2024-01-02T11:00:00Z"
}

variable "something" {
  type = string
}

# Set password_reset to be based on the update of another variable to reset the password
resource "ovh_cloud_project_database_clickhouse_user" "user_md5" {
  service_name    = data.ovh_cloud_project_database.clickhouse.service_name
  cluster_id      = data.ovh_cloud_project_database.clickhouse.id
  name            = "bob"
  password_reset  = md5(var.something)
}

# Change password_reset each time you want to reset the password to trigger an update
resource "ovh_cloud_project_database_clickhouse_user" "user" {
  service_name    = data.ovh_cloud_project_database.clickhouse.service_name
  cluster_id      = data.ovh_cloud_project_database.clickhouse.id
  name            = "johndoe"
  password_reset  = "reset1"
}

output "user_password" {
  value     = ovh_cloud_project_database_clickhouse_user.user.password
  sensitive = true
}
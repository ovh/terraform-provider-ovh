data "ovh_cloud_project_database" "m3db" {
  service_name  = "XXX"
  engine        = "m3db"
  id            = "ZZZ"
}

# Change password_reset with the datetime each time you want to reset the password to trigger an update
resource "ovh_cloud_project_database_m3db_user" "userDatetime" {
  service_name    = data.ovh_cloud_project_database.m3db.service_name
  cluster_id      = data.ovh_cloud_project_database.m3db.id
  group           = "mygroup"
  name            = "alice"
  password_reset  = "2024-01-02T11:00:00Z"
}

variable "something" {
  type = string
}

# Set password_reset to be based on the update of another variable to reset the password
resource "ovh_cloud_project_database_m3db_user" "userMd5" {
  service_name    = data.ovh_cloud_project_database.m3db.service_name
  cluster_id      = data.ovh_cloud_project_database.m3db.id
  group           = "mygroup"
  name            = "bob"
  password_reset  = md5(var.something)
}

# Change password_reset each time you want to reset the password to trigger an update
resource "ovh_cloud_project_database_m3db_user" "user" {
  service_name    = data.ovh_cloud_project_database.m3db.service_name
  cluster_id      = data.ovh_cloud_project_database.m3db.id
  group           = "mygroup"
  name            = "johndoe"
  password_reset  = "reset1"
}

output "user_password" {
    value = ovh_cloud_project_database_m3db_user.user.password
    sensitive = true
}

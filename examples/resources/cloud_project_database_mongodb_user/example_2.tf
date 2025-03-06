data "ovh_cloud_project_database" "mongodb" {
  service_name  = "XXX"
  engine        = "mongodb"
  id            = "ZZZ"
}

# Change password_reset with the datetime each time you want to reset the password to trigger an update
resource "ovh_cloud_project_database_mongodb_user" "userDatetime" {
  service_name  = data.ovh_cloud_project_database.mongodb.service_name
  cluster_id    = data.ovh_cloud_project_database.mongodb.id
  name          = "alice"
  roles         = ["backup@admin", "readAnyDatabase@admin"]
  password_reset  = "2024-01-02T11:00:00Z"
}

variable "something" {
  type = string
}

# Set password_reset to be based on the update of another variable to reset the password
resource "ovh_cloud_project_database_mongodb_user" "userMd5" {
  service_name  = data.ovh_cloud_project_database.mongodb.service_name
  cluster_id    = data.ovh_cloud_project_database.mongodb.id
  name          = "bob"
  roles         = ["backup@admin", "readAnyDatabase@admin"]
  password_reset  = md5(var.something)
}

# Change password_reset each time you want to reset the password to trigger an update
resource "ovh_cloud_project_database_mongodb_user" "user" {
  service_name  = data.ovh_cloud_project_database.mongodb.service_name
  cluster_id    = data.ovh_cloud_project_database.mongodb.id
  name          = "johndoe"
  roles         = ["backup@admin", "readAnyDatabase@admin"]
  password_reset  = "reset1"
}

output "user_password" {
  value     = ovh_cloud_project_database_mongodb_user.user.password
  sensitive = true
}

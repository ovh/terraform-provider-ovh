
data "ovh_cloud_project_database" "valkey" {
  service_name  = "XXXX"
  engine        = "valkey"
  id            = "ZZZZ"
}

# Change password_reset with the datetime each time you want to reset the password to trigger an update
resource "ovh_cloud_project_database_valkey_user" "user_datetime" {
  service_name    = data.ovh_cloud_project_database.valkey.service_name
  cluster_id      = data.ovh_cloud_project_database.valkey.id
  categories      = ["+@set", "+@sortedset"]
  channels        = ["*"]
  commands        = ["+get", "-set"]
  keys            = ["data", "properties"]
  name            = "alice"
  password_reset  = "2024-01-02T11:00:00Z"
}

variable "something" {
  type = string
}

# Set password_reset to be based on the update of another variable to reset the password
resource "ovh_cloud_project_database_valkey_user" "user_md5" {
  service_name    = data.ovh_cloud_project_database.valkey.service_name
  cluster_id      = data.ovh_cloud_project_database.valkey.id
  categories      = ["+@set", "+@sortedset"]
  channels        = ["*"]
  commands        = ["+get", "-set"]
  keys            = ["data", "properties"]
  name            = "bob"
  password_reset  = md5(var.something)
}

# Change password_reset each time you want to reset the password to trigger an update
resource "ovh_cloud_project_database_valkey_user" "user" {
  service_name    = data.ovh_cloud_project_database.valkey.service_name
  cluster_id      = data.ovh_cloud_project_database.valkey.id
  categories      = ["+@set", "+@sortedset"]
  channels        = ["*"]
  commands        = ["+get", "-set"]
  keys            = ["data", "properties"]
  name            = "johndoe"
  password_reset  = "reset1"
}

output "user_password" {
  value = ovh_cloud_project_database_valkey_user.user.password
  sensitive = true
}

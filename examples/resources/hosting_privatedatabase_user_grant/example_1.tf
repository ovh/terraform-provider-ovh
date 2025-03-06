resource "ovh_hosting_privatedatabase_user_grant" "user_grant" {
  service_name  = "XXXXXX"
  user_name     = "terraform"
  database_name = "ovhcloud"
  grant         = "admin"
}

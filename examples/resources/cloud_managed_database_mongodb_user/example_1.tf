resource "ovh_cloud_managed_database_mongodb_user" "app" {
  service_name   = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  service_id     = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  name           = "app"
  roles          = ["readAnyDatabase@admin"]
  password_reset = "2026-11-12T10:00:00Z"
}

output "app_password" {
  value     = ovh_cloud_managed_database_mongodb_user.app.password
  sensitive = true
}

# Change password_reset with the datetime each time you want to reset the password to trigger an update
resource "ovh_cloud_project_user" "user_with_rotation" {
  service_name   = "XXX"
  description    = "Service User created by Terraform with password rotation"
  password_reset = "2025-04-30"
}


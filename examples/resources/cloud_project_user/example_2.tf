resource "ovh_cloud_project_user" "user_with_rotation" {
  service_name = "XXX"
  description  = "Service User created by Terraform with password rotation"
  rotate_when_changed = {
    last_rotation = "2025-04-30"
  }
}


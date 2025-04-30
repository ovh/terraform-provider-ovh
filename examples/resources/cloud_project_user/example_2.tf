resource "ovh_cloud_project_user" "user_with_rotation" {
  service_name = "XXX"
  description  = "Service User created by Terraform with password rotation"

  # Rotate the password whenever the last_rotation value changes
  # This allows for scheduled rotation without recreating the user
  rotate_when_changed = {
    last_rotation = var.last_rotation
  }
}

# Variable to control password rotation
variable "last_rotation" {
  description = "Timestamp or other value that, when changed, will trigger password rotation"
  type        = string
  default     = "2025-04-30" # Initial rotation date
}

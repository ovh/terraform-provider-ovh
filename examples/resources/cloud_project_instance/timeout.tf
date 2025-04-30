resource "ovh_cloud_project_instance" "instance" {
  # ...

  timeouts {
    create = "10min"
  }
}
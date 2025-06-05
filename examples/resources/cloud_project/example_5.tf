resource "ovh_cloud_project" "my_cloud_project" {
  # ...

  timeouts {
    create = "1h"
  }
}
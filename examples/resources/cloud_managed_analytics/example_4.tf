resource "ovh_cloud_managed_analytics" "db" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

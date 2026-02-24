resource "ovh_cloud_managed_analytics_m3db_user" "user" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

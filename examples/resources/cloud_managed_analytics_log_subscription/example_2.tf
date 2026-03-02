resource "ovh_cloud_managed_analytics_log_subscription" "sub" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

resource "ovh_cloud_managed_analytics_ip_restriction" "ip_restriction" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

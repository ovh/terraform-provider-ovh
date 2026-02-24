resource "ovh_cloud_managed_analytics_m3db_namespace" "namespace" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

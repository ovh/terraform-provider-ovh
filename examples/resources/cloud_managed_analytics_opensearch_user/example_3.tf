resource "ovh_cloud_managed_analytics_opensearch_user" "user" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

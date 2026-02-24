resource "ovh_cloud_managed_analytics_opensearch_pattern" "pattern" {
  # ...

  timeouts {
    create = "1h"
    delete = "45m"
  }
}

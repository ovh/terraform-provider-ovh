resource "ovh_cloud_managed_analytics_kafka_acl" "acl" {
  # ...

  timeouts {
    create = "1h"
    delete = "45m"
  }
}

resource "ovh_cloud_managed_analytics_kafka_topic" "topic" {
  # ...

  timeouts {
    create = "1h"
    delete = "45m"
  }
}

resource "ovh_cloud_project_database_kafka_topic" "topic" {
  # ...

  timeouts {
    create = "1h"
    delete = "45m"
  }
}

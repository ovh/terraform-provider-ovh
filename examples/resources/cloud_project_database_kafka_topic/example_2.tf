resource "ovh_cloud_project_database_kafka_topic" "topic" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "45m"
  }
}

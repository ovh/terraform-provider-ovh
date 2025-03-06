resource "ovh_cloud_project_database_kafka_acl" "acl" {
  # ...

  timeouts {
    create = "1h"
    delete = "45m"
  }
}

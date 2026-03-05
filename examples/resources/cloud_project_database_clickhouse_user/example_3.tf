resource "ovh_cloud_project_database_clickhouse_user" "user" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

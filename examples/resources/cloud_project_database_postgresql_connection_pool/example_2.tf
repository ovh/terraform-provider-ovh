resource "ovh_cloud_project_database_postgresql_connection_pool" "user" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

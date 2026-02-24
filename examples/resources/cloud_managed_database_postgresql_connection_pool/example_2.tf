resource "ovh_cloud_managed_database_postgresql_connection_pool" "user" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

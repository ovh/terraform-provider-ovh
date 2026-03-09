resource "ovh_cloud_managed_database_postgresql_user" "user" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

resource "ovh_cloud_managed_database_database" "database" {
  # ...

  timeouts {
    create = "1h"
    delete = "45m"
  }
}

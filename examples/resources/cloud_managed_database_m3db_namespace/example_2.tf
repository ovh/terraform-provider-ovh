resource "ovh_cloud_managed_database_m3db_namespace" "namespace" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

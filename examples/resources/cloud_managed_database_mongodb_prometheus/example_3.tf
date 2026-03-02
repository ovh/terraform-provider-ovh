resource "ovh_cloud_managed_database_mongodb_prometheus" "prometheus" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

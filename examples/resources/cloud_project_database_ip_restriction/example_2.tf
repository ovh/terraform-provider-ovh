resource "ovh_cloud_project_database_ip_restriction" "ip_restriction" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

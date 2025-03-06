resource "ovh_cloud_project_database_opensearch_pattern" "pattern" {
  # ...

  timeouts {
    create = "1h"
    delete = "45m"
  }
}

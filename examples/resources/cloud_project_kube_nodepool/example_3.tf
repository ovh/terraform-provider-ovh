resource "ovh_cloud_project_kube_nodepool" "pool" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

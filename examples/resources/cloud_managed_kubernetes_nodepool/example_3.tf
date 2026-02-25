resource "ovh_cloud_managed_kubernetes_nodepool" "pool" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

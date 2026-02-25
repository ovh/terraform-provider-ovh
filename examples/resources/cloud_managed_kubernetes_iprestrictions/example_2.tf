resource "ovh_cloud_managed_kubernetes_iprestrictions" "vrack_only" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

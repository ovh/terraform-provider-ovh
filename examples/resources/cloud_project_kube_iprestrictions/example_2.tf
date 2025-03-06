resource "ovh_cloud_project_kube_iprestrictions" "vrack_only" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

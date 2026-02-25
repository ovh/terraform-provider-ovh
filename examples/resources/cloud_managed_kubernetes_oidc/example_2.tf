resource "ovh_cloud_managed_kubernetes_oidc" "oidc" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

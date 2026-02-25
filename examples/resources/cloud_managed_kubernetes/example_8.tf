resource "ovh_cloud_managed_kubernetes" "my_kube_cluster" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

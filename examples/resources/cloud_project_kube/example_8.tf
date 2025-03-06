resource "ovh_cloud_project_kube" "my_kube_cluster" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

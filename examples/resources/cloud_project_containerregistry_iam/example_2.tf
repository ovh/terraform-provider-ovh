resource "ovh_cloud_project_containerregistry_iam" "registry_iam" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

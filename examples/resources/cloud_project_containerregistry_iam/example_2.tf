resource "ovh_cloud_project_containerregistry_iam" "my-iam" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

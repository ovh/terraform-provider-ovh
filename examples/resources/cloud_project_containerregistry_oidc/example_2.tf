resource "ovh_cloud_project_containerregistry_oidc" "my_oidc" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}

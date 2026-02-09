resource "ovh_storage_efs" "efs" {
  # ...

  timeouts {
    create = "1h"
  }
}

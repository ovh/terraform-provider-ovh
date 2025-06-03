resource "ovh_vrack" "vrack" {
  # ...

  timeouts {
    create = "1h"
  }
}
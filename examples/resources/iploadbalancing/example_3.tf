resource "ovh_iploadbalancing" "iplb" {
  # ...

  timeouts {
    create = "1h"
  }
}
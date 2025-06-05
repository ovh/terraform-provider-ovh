resource "ovh_domain_zone" "zone" {
  # ...

  timeouts {
    create = "1h"
  }
}
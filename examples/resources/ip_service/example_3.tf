resource "ovh_ip_service" "ipblock" {
  # ...

  timeouts {
    create = "1h"
  }
}
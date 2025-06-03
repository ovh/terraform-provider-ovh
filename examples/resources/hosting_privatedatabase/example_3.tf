resource "ovh_hosting_privatedatabase" "database" {
  # ...

  timeouts {
    create = "1h"
  }
}
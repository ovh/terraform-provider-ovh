resource "ovh_vps_set_password" "pwd" {
  service_name = "vpsXXXXXX.vps.ovh.net"

  # Re-run the action by changing any value in this map.
  triggers = {
    rotate = "2025-01-01"
  }
}

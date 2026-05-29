resource "ovh_vps_reboot" "reboot" {
  service_name = "vpsXXXXXX.vps.ovh.net"

  # Change any value here to re-run the reboot.
  triggers = {
    kernel = "6.6.1"
  }
}

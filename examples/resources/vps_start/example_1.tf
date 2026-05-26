resource "ovh_vps_start" "start" {
  service_name = "vpsXXXXXX.vps.ovh.net"

  triggers = {
    nonce = "1"
  }
}

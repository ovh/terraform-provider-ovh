resource "ovh_vps_stop" "stop" {
  service_name = "vpsXXXXXX.vps.ovh.net"

  triggers = {
    nonce = "1"
  }
}

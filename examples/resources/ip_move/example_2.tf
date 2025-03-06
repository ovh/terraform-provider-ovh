resource "ovh_ip_move" "park_ip" {
  ip = "1.2.3.4"
  routed_to {
    service_name = ""
  }
}

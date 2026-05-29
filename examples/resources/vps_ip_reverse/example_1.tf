resource "ovh_vps_ip_reverse" "rev" {
  service_name = "vpsXXXXX.ovh.net"
  ip_address   = "192.0.2.1"
  reverse      = "host.example.com."
}

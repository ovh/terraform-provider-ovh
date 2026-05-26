resource "ovh_vps_secondary_dns_domain" "example" {
  service_name = "vpsXXXXX.ovh.net"
  domain       = "example.com"
  ip           = "203.0.113.10"
}

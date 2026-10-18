resource "ovh_ip_firewall_rule" "allow_app_ports" {
  ip             = "XXXXXX"
  ip_on_firewall = "XXXXXX"
  sequence       = 1
  action         = "permit"
  protocol       = "tcp"

  destination_port_range = {
    from = 8000
    to   = 8010
  }
}

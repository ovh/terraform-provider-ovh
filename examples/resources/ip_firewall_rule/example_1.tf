resource "ovh_ip_firewall_rule" "my_firewall_rule" {
  ip             = "XXXXXX"
  ip_on_firewall = "XXXXXX"
  sequence       = 0
  action         = "deny"
  protocol       = "tcp"
}

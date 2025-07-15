resource "ovh_ip_firewall_rule" "deny_tcp" {
  ip             = "XXXXXX"
  ip_on_firewall = "XXXXXX"
  sequence       = 0
  action         = "deny"
  protocol       = "tcp"
}

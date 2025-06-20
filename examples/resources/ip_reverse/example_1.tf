# Set the reverse of an IP
resource "ovh_ip_reverse" "test" {
  readiness_timeout_duration = "1m"
  ip                         = "192.0.2.0/24"
  ip_reverse                 = "192.0.2.1"
  reverse                    = "example.com"
}

resource "ovh_domain_zone_dynhost_record" "dynhost_record" {
  zone_name  = "mydomain.ovh"
  sub_domain = "dynhost"
  ip         = "1.2.3.4"
}

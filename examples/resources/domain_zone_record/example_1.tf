# Add a record to a sub-domain
resource "ovh_domain_zone_record" "test" {
  zone      = "testdemo.ovh"
  subdomain = "test"
  fieldtype = "A"
  ttl       = 3600
  target    = "0.0.0.0"
}

# Add a redirection to a sub-domain
resource "ovh_domain_zone_redirection" "test" {
  zone      = "testdemo.ovh"
  subdomain = "test"
  type      = "visiblePermanent"
  target    = "http://www.ovh"
}

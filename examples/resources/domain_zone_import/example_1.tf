resource "ovh_domain_zone_import" "import" {
  zone_name = "mysite.ovh"
  zone_file = file("./example.zone")
}

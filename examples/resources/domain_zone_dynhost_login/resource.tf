resource "ovh_domain_zone_dynhost_login" "dynhost_user" {
  zone_name    = "mydomain.ovh"
  sub_domain   = "dynhost"
  login_suffix = "dynhostUser"
  password     = "thisIsMyPassword"
}

data "ovh_vps_templates" "all" {
  service_name = "vps-xxxxxx.vps.ovh.net"
}

data "ovh_vps_templates" "debian_64" {
  service_name        = "vps-xxxxxx.vps.ovh.net"
  distribution_filter = "Debian"
  bit_format_filter   = 64
}

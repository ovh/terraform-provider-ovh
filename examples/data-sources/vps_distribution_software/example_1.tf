data "ovh_vps_distribution_software" "installed" {
  service_name  = "vpsXXXXX.ovh.net"
  type_filter   = "webserver"
  status_filter = "stable"
}

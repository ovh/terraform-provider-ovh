data "ovh_vps_service_info" "info" {
  service_name = "vpsXXXXX.ovh.net"
}

output "expiration_date" {
  value = data.ovh_vps_service_info.info.expiration
}

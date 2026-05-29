data "ovh_vrack_vpss" "vpss" {
  service_name = "pn-XXXXXX"
}

output "vps_list" {
  value = data.ovh_vrack_vpss.vpss.vpss
}

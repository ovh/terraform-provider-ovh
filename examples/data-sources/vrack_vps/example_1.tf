data "ovh_vrack_vps" "vrack_vps" {
  service_name = "pn-XXXXXX"
  vps          = "vpsXXXXX.ovh.net"
}

output "vrack_vps_state" {
  value = data.ovh_vrack_vps.vrack_vps.state
}

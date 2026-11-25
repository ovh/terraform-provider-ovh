data "ovh_vps" "my_vps" {
  id = "vps-1a2b3c4d.vps.ovh.net"
}

output "vps_region" {
  value = data.ovh_vps.my_vps.current_state.location.region
}

output "vps_ipv4" {
  value = [for address in data.ovh_vps.my_vps.current_state.ip_addresses : address.ip if address.version == 4]
}

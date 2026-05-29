data "ovh_vps_snapshot" "current" {
  service_name = "vpsXXXXX.ovh.net"
}

output "snapshot_description" {
  value = data.ovh_vps_snapshot.current.description
}

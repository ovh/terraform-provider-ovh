data "ovh_vps_automated_backup" "backup" {
  service_name = "vpsXXXXX.ovh.net"
}

output "backup_state" {
  value = data.ovh_vps_automated_backup.backup.state
}

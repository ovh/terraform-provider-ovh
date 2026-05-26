data "ovh_vps_automated_backup_attached" "attached" {
  service_name = "vpsXXXXX.ovh.net"
}

output "attached_restore_point" {
  value = data.ovh_vps_automated_backup_attached.attached.restore_point
}

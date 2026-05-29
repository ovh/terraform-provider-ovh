resource "ovh_vps_automated_backup_restore" "restore" {
  service_name  = "vpsXXXXX.ovh.net"
  restore_point = "2024-01-15T02:00:00+00:00"
  type          = "file"
  change_password = false
}

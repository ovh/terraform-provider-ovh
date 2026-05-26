resource "ovh_vps_automated_backup_reschedule" "schedule" {
  service_name = "vpsXXXXXX.vps.ovh.net"
  schedule     = "02:00:00"
}

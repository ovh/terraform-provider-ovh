data "ovh_vps_automated_backup_restore_points" "restore_points" {
  service_name = "vpsXXXXX.ovh.net"
}

output "restore_point_ids" {
  value = data.ovh_vps_automated_backup_restore_points.restore_points.restore_points
}

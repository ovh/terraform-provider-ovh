data "ovh_vps_backup_ftp" "ftp" {
  service_name = "vpsXXXXX.ovh.net"
}

output "ftp_usage" {
  value = data.ovh_vps_backup_ftp.ftp.usage
}

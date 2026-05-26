resource "ovh_vps_backup_ftp_password" "rotate" {
  service_name = "vpsXXXXXX.vps.ovh.net"

  triggers = {
    rotation = "1"
  }
}

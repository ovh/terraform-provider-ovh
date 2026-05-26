resource "ovh_vps_backup_ftp_access" "access" {
  service_name = "vpsXXXXX.ovh.net"
  ip_block     = "198.51.100.0/24"
  ftp          = true
  nfs          = false
  cifs         = false
}

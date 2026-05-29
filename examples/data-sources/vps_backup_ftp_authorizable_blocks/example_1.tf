data "ovh_vps_backup_ftp_authorizable_blocks" "blocks" {
  service_name = "vpsXXXXX.ovh.net"
}

output "authorizable_blocks" {
  value = data.ovh_vps_backup_ftp_authorizable_blocks.blocks.blocks
}

data "ovh_vps_snapshot_download" "download" {
  service_name = "vpsXXXXX.ovh.net"
}

output "download_url" {
  value     = data.ovh_vps_snapshot_download.download.url
  sensitive = true
}

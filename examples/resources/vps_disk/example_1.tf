resource "ovh_vps_disk" "disk" {
  service_name             = "vps-XXXXXX.vps.ovh.net"
  disk_id                  = 1234
  monitoring               = true
  low_free_space_threshold = 1024
}

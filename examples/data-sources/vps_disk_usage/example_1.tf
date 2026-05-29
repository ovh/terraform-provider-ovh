data "ovh_vps_disk_usage" "usage" {
  service_name = "vps-XXXXXX.vps.ovh.net"
  disk_id      = 1234
  type         = "used"
}

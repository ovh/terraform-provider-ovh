data "ovh_vps_disk_monitoring" "mon" {
  service_name = "vps-XXXXXX.vps.ovh.net"
  disk_id      = 1234
  period       = "lastday"
  type         = "cpu:used"
}

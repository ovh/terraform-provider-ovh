data "ovh_vps_tasks" "tasks" {
  service_name = "vps-XXXXXX.vps.ovh.net"
  state_filter = "doing"
  type_filter  = "rebootVm"
}

data "ovh_dedicated_server_boots" "rescue" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
  boot_type    = "rescue"
  kernel       = "rescue64-pro"
}

resource "ovh_dedicated_server_update" "server_on_rescue" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
  boot_id      = data.ovh_dedicated_server_boots.rescue.result[0]
  monitoring   = true
  state        = "ok"
}

resource "ovh_dedicated_server_reboot_task" "server_reboot" {
  service_name = data.ovh_dedicated_server_boots.rescue.service_name

  keepers = [
     ovh_dedicated_server_update.server_on_rescue.boot_id,
  ]
}

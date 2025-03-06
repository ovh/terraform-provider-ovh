data "ovh_dedicated_server_boots" "rescue" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
  boot_type    = "rescue"
  kernel       = "rescue64-pro"
}

resource "ovh_dedicated_server_update" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
  boot_id      = data.ovh_dedicated_server_boots.rescue.result[0]
  monitoring   = true
  state        = "ok"
  display_name = "Some human-readable name"
}

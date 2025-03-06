data "ovh_dedicated_server" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

resource "ovh_vrack_dedicated_server_interface" "vdsi" {
  service_name = "pn-xxxxxxx" #name of the vrack
  interface_id = data.ovh_dedicated_server.server.enabled_vrack_vnis[0]
}

data "ovh_dedicated_server" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

resource "ovh_dedicated_server_networking" "server" {
  service_name = data.ovh_dedicated_server.server.service_name

  interfaces {
    macs = slice(sort(flatten(data.ovh_dedicated_server.server.vnis.*.nics)), 0, 2)
    type = "public"
  }
  interfaces {
    macs = slice(sort(flatten(data.ovh_dedicated_server.server.vnis.*.nics)), 2, 4)
    type = "vrack"
  }
}

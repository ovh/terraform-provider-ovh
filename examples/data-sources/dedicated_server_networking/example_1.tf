data "ovh_dedicated_server_networking" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

output "public_lacp_fallback_mac" {
  value = one([
    for iface in data.ovh_dedicated_server_networking.server.interfaces :
    iface.aggregation_fallback if iface.type == "public"
  ])
}

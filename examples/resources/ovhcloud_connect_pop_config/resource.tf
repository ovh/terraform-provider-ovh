data "ovh_ovhcloud_connect" "occ" {
   service_name = "XXX"
}

resource "ovh_ovhcloud_connect_pop_config" "pop" {
    service_name = data.ovh_ovhcloud_connect.occ.service_name
    interface_id = tolist(data.ovh_ovhcloud_connect.occ.interface_list)[0]
    type = "l3"
    customer_bgp_area = 65400
    ovh_bgp_area = 65401
    subnet = "I.I.I.I/30"
}

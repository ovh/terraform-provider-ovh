data "ovh_ovhcloud_connect_config_pops" "pop_cfgs" {
  service_name = "XXX"
}

data "ovh_ovhcloud_connect_config_pop_datacenters" "datacenter_cfgs" {
  service_name = data.ovh_ovhcloud_connect_config_pops.pop_cfgs.service_name
  config_pop_id = tolist(data.ovh_ovhcloud_connect_config_pops.pop_cfgs.pop_configs)[0].id
}

resource "ovh_ovhcloud_connect_pop_datacenter_extra_config" "extra" {
    service_name = data.ovh_ovhcloud_connect_config_pops.pop_cfgs.service_name
    config_pop_id = tolist(data.ovh_ovhcloud_connect_config_pops.pop_cfgs.pop_configs)[0].id
    config_datacenter_id = tolist(data.ovh_ovhcloud_connect_config_pop_datacenters.datacenter_cfgs.datacenter_configs)[0].id
    type = "network"
    next_hop = "P.P.P.P"
    subnet = "I.I.I.I/M"
}

data "ovh_ovhcloud_connect_config_pops" "pop_cfgs" {
  service_name = "XXX"
}

data "ovh_ovhcloud_connect_datacenters" "dcs" {
  service_name = data.ovh_ovhcloud_connect_config_pops.pop_cfgs.service_name
}

output "occ_datacenters" {
  value = data.ovh_ovhcloud_connect_datacenters.dcs
}

resource "ovh_ovhcloud_connect_pop_datacenter_config" "dc" {
  service_name = data.ovh_ovhcloud_connect_config_pops.pop_cfgs.service_name
  config_pop_id = tolist(data.ovh_ovhcloud_connect_config_pops.pop_cfgs.pop_configs)[0].id
  datacenter_id = tolist(data.ovh_ovhcloud_connect_datacenters.dcs)[0].id
  ovh_bgp_area = 65408
  subnet = "I.I.I.I/28"
}

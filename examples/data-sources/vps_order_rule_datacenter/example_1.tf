data "ovh_vps_order_rule_datacenter" "dc_rule" {
  plan_code = "vps-starter-1-2-20"
  ovh_subsidiary = "FR"
}

output "available_datacenters" {
  value = data.ovh_vps_order_rule_datacenter.dc_rule.datacenters
}

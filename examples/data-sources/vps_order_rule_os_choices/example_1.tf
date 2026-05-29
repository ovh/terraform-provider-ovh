data "ovh_vps_order_rule_os_choices" "os_choices" {
  plan_code      = "vps-starter-1-2-20"
  ovh_subsidiary = "FR"
  datacenter     = "GRA"
}

output "os_image_ids" {
  value = data.ovh_vps_order_rule_os_choices.os_choices.images[*].id
}

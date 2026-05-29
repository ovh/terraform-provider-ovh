data "ovh_vps_available_upgrade" "upgrade" {
  service_name = "vpsXXXXX.ovh.net"
}

output "upgrade_offers" {
  value = data.ovh_vps_available_upgrade.upgrade.plans[*].plan_code
}

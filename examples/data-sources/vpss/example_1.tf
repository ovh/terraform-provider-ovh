data "ovh_vpss" "all" {}

output "vps_ids" {
  value = data.ovh_vpss.all.vps[*].id
}

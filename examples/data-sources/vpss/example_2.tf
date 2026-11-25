data "ovh_vpss" "all" {}

output "vps_in_gravelines" {
  value = [for vps in data.ovh_vpss.all.vps : vps.id if vps.location.region == "eu-west-gra"]
}

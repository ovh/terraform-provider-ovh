data "ovh_vps_datacenters" "dc" {
  service_name = "vpsXXXXX.ovh.net"
}

output "datacenter_names" {
  value = data.ovh_vps_datacenters.dc.datacenters[*].name
}

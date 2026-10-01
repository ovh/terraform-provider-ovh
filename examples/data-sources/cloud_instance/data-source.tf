data "ovh_cloud_instance" "instance" {
  service_name = "<Public cloud project id>"
  id           = "<instance id>"
}

output "instance_flavor_name" {
  value = data.ovh_cloud_instance.instance.current_state.flavor.name
}

# Every address observed on the instance, all interfaces and address types
# (FIXED, FLOATING, ADDITIONAL) combined.
output "instance_addresses" {
  value = flatten([
    for network in data.ovh_cloud_instance.instance.current_state.networks :
    [for address in network.addresses : address.ip]
  ])
}

# Unlike ovh_cloud_instances, the singular data source populates
# current_state.shares.
output "instance_share_ids" {
  value = [
    for share in data.ovh_cloud_instance.instance.current_state.shares :
    share.id
  ]
}

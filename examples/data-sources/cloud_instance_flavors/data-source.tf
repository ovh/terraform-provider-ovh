# Omit `region` to list the catalog of every region at once: the same flavor
# then appears once per region it is offered in.
data "ovh_cloud_instance_flavors" "flavors" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
}

# Resolve a flavor id by its commercial name, to feed ovh_cloud_instance.flavor_id.
output "b3_8_flavor_id" {
  value = one([
    for flavor in data.ovh_cloud_instance_flavors.flavors.flavors :
    flavor.id if flavor.name == "b3-8"
  ])
}

# Every public flavor of the region with at least 16 GB of RAM.
output "large_flavor_names" {
  value = sort([
    for flavor in data.ovh_cloud_instance_flavors.flavors.flavors :
    flavor.name if flavor.is_public && flavor.ram >= 16384
  ])
}

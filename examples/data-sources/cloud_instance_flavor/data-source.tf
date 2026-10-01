data "ovh_cloud_instance_flavor" "flavor" {
  service_name = "<Public cloud project id>"
  id           = "<flavor id>"
}

output "flavor_sizing" {
  value = format(
    "%s: %d vCPUs, %d MB RAM, %d GB disk",
    data.ovh_cloud_instance_flavor.flavor.name,
    data.ovh_cloud_instance_flavor.flavor.vcpus,
    data.ovh_cloud_instance_flavor.flavor.ram,
    data.ovh_cloud_instance_flavor.flavor.disk,
  )
}

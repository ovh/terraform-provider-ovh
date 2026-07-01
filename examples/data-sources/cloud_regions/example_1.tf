data "ovh_cloud_regions" "regions" {
  service_name = "<public cloud project ID>"
}

output "regions" {
  value = data.ovh_cloud_regions.regions.regions
}

data "ovh_cloud_project_loadbalancers" "lbs" {
  service_name = "XXXXXX"
  region_name  = "XXX"
}
output "lbs" {
  value = data.ovh_cloud_project_loadbalancers.lbs
}

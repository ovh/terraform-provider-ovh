data "ovh_cloud_project_loadbalancer" "lb" {
  service_name = "XXXXXX"
  region_name  = "XXX"
  id           = "XXX"
}
output "lb" {
  value = data.ovh_cloud_project_loadbalancer.lb
}

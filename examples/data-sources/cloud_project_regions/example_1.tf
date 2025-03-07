data "ovh_cloud_project_regions" "regions" {
  service_name    = "XXXXXX"
  has_services_up = ["network"]
}

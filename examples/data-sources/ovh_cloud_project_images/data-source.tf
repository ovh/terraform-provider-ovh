data "ovh_cloud_project_images" "images" {
  service_name = "<public cloud project ID>"
  region       = "WAW1"
  os_type      = "linux"
}
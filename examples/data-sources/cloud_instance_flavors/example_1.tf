data "ovh_cloud_instance_flavors" "b2" {
  service_name = "YYYY"
  region_name  = "GRA7"
  name         = "b2-.*"
}

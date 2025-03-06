data "ovh_cloud_project_storage_object" "object" {
  service_name = "<public cloud project ID>"
  region_name  = "GRA"
  name         = "<bucket name>"
  key          = "<object name>"
}
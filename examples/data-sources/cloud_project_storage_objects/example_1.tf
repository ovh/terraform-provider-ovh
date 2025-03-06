data "ovh_cloud_project_storage_objects" "objects" {
  service_name = "<public cloud project ID>"
  region_name  = "GRA"
  name         = "<bucket name>"
}

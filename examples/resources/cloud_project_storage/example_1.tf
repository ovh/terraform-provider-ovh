resource "ovh_cloud_project_storage" "storage" {
  service_name = "<public cloud project ID>"
  region_name = "GRA"
  name = "my-storage"
  versioning = {
    status = "enabled"
  }
}

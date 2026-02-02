resource "ovh_cloud_project_storage" "storage_with_lock" {
  service_name = "<public cloud project ID>"
  region_name  = "GRA"
  name         = "my-encrypted-storage"

  encryption = {
    sse_algorithm = "AES256"
  }
}

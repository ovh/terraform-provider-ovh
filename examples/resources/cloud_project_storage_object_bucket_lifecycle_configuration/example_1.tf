resource "ovh_cloud_project_storage" "bucket" {
  service_name = "<public cloud project ID>"
  region_name  = "GRA"
  name         = "my-versioned-bucket"
  versioning = {
    status = "enabled"
  }
}

resource "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
  service_name   = ovh_cloud_project_storage.bucket.service_name
  region_name    = ovh_cloud_project_storage.bucket.region_name
  container_name = ovh_cloud_project_storage.bucket.name

  rules = [
    {
      id     = "expire-old-objects"
      status = "enabled"
      expiration = {
        days = 365
      }
    }
  ]
}
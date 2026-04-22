data "ovh_cloud_storage_object_bucket_objects" "example" {
  service_name = "<public cloud project ID>"
  bucket_name  = "my-bucket"

  prefix    = "logs/"
  delimiter = "/"
  limit     = 100
}

# List all versions (and delete markers):
data "ovh_cloud_storage_object_bucket_objects" "versions" {
  service_name  = "<public cloud project ID>"
  bucket_name   = "my-bucket"
  with_versions = true
}

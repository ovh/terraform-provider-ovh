data "ovh_cloud_storage_object_bucket_object" "example" {
  service_name = "<public cloud project ID>"
  bucket_name  = "my-bucket"
  object_key   = "path/to/my-object.txt"
}

# Targeting a specific version:
data "ovh_cloud_storage_object_bucket_object" "example_version" {
  service_name = "<public cloud project ID>"
  bucket_name  = "my-bucket"
  object_key   = "path/to/my-object.txt"
  version_id   = "1a2b3c4d5e6f"
}

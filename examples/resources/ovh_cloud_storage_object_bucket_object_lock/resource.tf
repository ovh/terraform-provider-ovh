resource "ovh_cloud_storage_object_bucket_object_lock" "example" {
  service_name = "<public cloud project ID>"
  bucket_name  = "my-locked-bucket"
  object_key   = "path/to/my-object.txt"

  retention = {
    mode              = "GOVERNANCE"
    retain_until_date = "2030-01-01T00:00:00Z"
  }

  legal_hold = {
    status = "ON"
  }
}

# Targeting a specific version:
resource "ovh_cloud_storage_object_bucket_object_lock" "example_version" {
  service_name = "<public cloud project ID>"
  bucket_name  = "my-locked-bucket"
  object_key   = "path/to/my-object.txt"
  version_id   = "1a2b3c4d5e6f"

  legal_hold = {
    status = "ON"
  }
}

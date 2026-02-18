resource "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
  service_name   = "<public cloud project ID>"
  region_name    = "GRA"
  container_name = "my-bucket"

  rules = [
    {
      id     = "noncurrent-version-expiration"
      status = "enabled"
      noncurrent_version_expiration = {
        noncurrent_days           = 30
        newer_noncurrent_versions = 3
      }
    },
    {
      id     = "noncurrent-version-transition"
      status = "enabled"
      noncurrent_version_transitions = [
        {
          noncurrent_days = 14
          storage_class   = "STANDARD_IA"
        }
      ]
    }
  ]
}

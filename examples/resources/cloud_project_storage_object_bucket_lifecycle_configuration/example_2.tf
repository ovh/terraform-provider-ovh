resource "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
  service_name   = "<public cloud project ID>"
  region_name    = "GRA"
  container_name = "my-bucket"

  rules = [
    {
      id     = "transition-to-standard-ia"
      status = "enabled"
      filter = {
        prefix = "archive/"
      }
      transitions = [
        {
          days          = 90
          storage_class = "STANDARD_IA"
        }
      ]
    },
    {
      id     = "transition-to-standard-ia-delayed"
      status = "enabled"
      filter = {
        prefix = "archive/"
      }
      transitions = [
        {
          days          = 180
          storage_class = "DEEP_ARCHIVE"
        }
      ]
    },
    {
      id     = "expire-logs"
      status = "enabled"
      filter = {
        prefix                   = "logs/"
        object_size_greater_than = 1048576
      }
      expiration = {
        days = 30
      }
    },
    {
      id     = "abort-multipart"
      status = "enabled"
      abort_incomplete_multipart_upload = {
        days_after_initiation = 7
      }
    }
  ]
}

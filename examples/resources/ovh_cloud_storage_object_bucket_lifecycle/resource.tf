resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name = "<Public cloud project id>"
  name         = "my-versioned-bucket"
  region       = "GRA"

  versioning = {
    status = "ENABLED"
  }
}

resource "ovh_cloud_storage_object_bucket_lifecycle" "lifecycle" {
  service_name = ovh_cloud_storage_object_bucket.bucket.service_name
  bucket_id    = ovh_cloud_storage_object_bucket.bucket.id

  rules = [
    {
      id     = "expire-logs"
      status = "ENABLED"
      filter = {
        prefix = "logs/"
      }
      expiration = {
        days = 30
      }
    },
    {
      id     = "archive-reports"
      status = "ENABLED"
      filter = {
        prefix = "reports/"
        tags = {
          retention = "long"
        }
      }
      transitions = [
        {
          days          = 90
          storage_class = "STANDARD_IA"
        },
        {
          days          = 365
          storage_class = "DEEP_ARCHIVE"
        }
      ]
    },
    {
      id     = "clean-noncurrent-versions"
      status = "ENABLED"
      noncurrent_version_transitions = [
        {
          noncurrent_days = 14
          storage_class   = "STANDARD_IA"
        }
      ]
      noncurrent_version_expiration = {
        noncurrent_days           = 60
        newer_noncurrent_versions = 3
      }
      expiration = {
        expired_object_delete_marker = true
      }
    },
    {
      id     = "abort-multipart"
      status = "ENABLED"
      abort_incomplete_multipart_upload = {
        days_after_initiation = 7
      }
    }
  ]
}

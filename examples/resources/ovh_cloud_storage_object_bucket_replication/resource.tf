resource "ovh_cloud_storage_object_bucket" "replica" {
  service_name = "<Public cloud project id>"
  name         = "my-replica-bucket"
  region       = "SBG"

  versioning = {
    status = "ENABLED"
  }
}

resource "ovh_cloud_storage_object_bucket" "source" {
  service_name = "<Public cloud project id>"
  name         = "my-source-bucket"
  region       = "GRA"

  versioning = {
    status = "ENABLED"
  }
}

resource "ovh_cloud_storage_object_bucket_replication" "replication" {
  service_name = ovh_cloud_storage_object_bucket.source.service_name
  bucket_id    = ovh_cloud_storage_object_bucket.source.id

  rules = [
    {
      id                        = "replicate-all"
      status                    = "ENABLED"
      priority                  = 1
      delete_marker_replication = "DISABLED"
      destination = {
        name   = ovh_cloud_storage_object_bucket.replica.name
        region = ovh_cloud_storage_object_bucket.replica.region
      }
    },
    {
      id                        = "archive-logs"
      status                    = "ENABLED"
      priority                  = 2
      delete_marker_replication = "DISABLED"
      filter = {
        prefix = "logs/"
        tags = {
          retention = "long"
        }
      }
      destination = {
        name          = ovh_cloud_storage_object_bucket.replica.name
        region        = ovh_cloud_storage_object_bucket.replica.region
        storage_class = "STANDARD_IA"
      }
    }
  ]
}

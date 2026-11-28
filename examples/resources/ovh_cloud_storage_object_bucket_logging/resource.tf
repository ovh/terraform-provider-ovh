resource "ovh_cloud_storage_object_bucket" "logs" {
  service_name = "<Public cloud project id>"
  name         = "my-logs-bucket"
  region       = "GRA"
}

resource "ovh_cloud_storage_object_bucket_acl" "logs" {
  service_name = ovh_cloud_storage_object_bucket.logs.service_name
  bucket_id    = ovh_cloud_storage_object_bucket.logs.id
  acl          = "log-delivery-write"
}

resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name = "<Public cloud project id>"
  name         = "my-data-bucket"
  region       = "GRA"
}

resource "ovh_cloud_storage_object_bucket_logging" "logging" {
  service_name  = ovh_cloud_storage_object_bucket.bucket.service_name
  bucket_id     = ovh_cloud_storage_object_bucket.bucket.id
  target_bucket = ovh_cloud_storage_object_bucket.logs.name
  target_prefix = "my-data-bucket/"

  depends_on = [ovh_cloud_storage_object_bucket_acl.logs]
}

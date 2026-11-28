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

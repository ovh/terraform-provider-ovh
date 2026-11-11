resource "ovh_cloud_storage_object_bucket_replication_job" "catch_up" {
  service_name = ovh_cloud_storage_object_bucket_replication.replication.service_name
  bucket_id    = ovh_cloud_storage_object_bucket_replication.replication.bucket_id
}

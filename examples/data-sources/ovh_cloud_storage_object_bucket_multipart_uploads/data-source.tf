data "ovh_cloud_storage_object_bucket_multipart_uploads" "uploads" {
  service_name = "<Public cloud project id>"
  bucket_id    = "GRA_my-bucket"
}

output "multipart_upload_keys" {
  value = { for upload in data.ovh_cloud_storage_object_bucket_multipart_uploads.uploads.uploads : upload.id => upload.key }
}

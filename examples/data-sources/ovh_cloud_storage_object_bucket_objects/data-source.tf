data "ovh_cloud_storage_object_bucket_objects" "documents" {
  service_name = "<Public cloud project id>"
  bucket_id    = "GRA_my-bucket"
  prefix       = "documents/"
  delimiter    = "/"
}

output "document_keys" {
  value = [
    for object in data.ovh_cloud_storage_object_bucket_objects.documents.objects :
    object.key if !object.is_common_prefix
  ]
}

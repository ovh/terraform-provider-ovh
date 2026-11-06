data "ovh_cloud_storage_object_bucket_object" "report" {
  service_name = "<Public cloud project id>"
  bucket_id    = "GRA_my-bucket"
  key          = "documents/report.pdf"
}

output "report_size" {
  value = data.ovh_cloud_storage_object_bucket_object.report.size
}

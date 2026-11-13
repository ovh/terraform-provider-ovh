data "ovh_cloud_storage_object_bucket_replication" "replication" {
  service_name = "<Public cloud project id>"
  bucket_id    = "GRA_my-source-bucket"
}

output "replication_destinations" {
  value = [for rule in data.ovh_cloud_storage_object_bucket_replication.replication.rules : rule.destination.name]
}

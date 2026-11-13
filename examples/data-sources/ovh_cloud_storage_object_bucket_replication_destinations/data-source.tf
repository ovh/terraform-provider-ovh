data "ovh_cloud_storage_object_bucket_replication_destinations" "destinations" {
  service_name = "<Public cloud project id>"
  bucket_id    = "GRA_my-source-bucket"
}

output "eligible_destinations" {
  value = [
    for destination in data.ovh_cloud_storage_object_bucket_replication_destinations.destinations.destinations :
    "${destination.region}/${destination.name}" if destination.eligible
  ]
}

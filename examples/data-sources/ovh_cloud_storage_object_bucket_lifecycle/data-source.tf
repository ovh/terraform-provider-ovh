data "ovh_cloud_storage_object_bucket_lifecycle" "lifecycle" {
  service_name = "<Public cloud project id>"
  bucket_id    = "GRA_my-bucket"
}

output "lifecycle_rule_ids" {
  value = [for rule in data.ovh_cloud_storage_object_bucket_lifecycle.lifecycle.rules : rule.id]
}

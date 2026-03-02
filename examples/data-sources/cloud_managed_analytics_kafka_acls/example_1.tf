data "ovh_cloud_managed_analytics_kafka_acls" "acls" {
  service_name  = "XXX"
  cluster_id    = "YYY"
}

output "acl_ids" {
  value = data.ovh_cloud_managed_analytics_kafka_acls.acls.acl_ids
}

data "ovh_cloud_managed_analytics_kafka_acl" "acl" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  id            = "ZZZ"
}

output "acl_permission" {
  value = data.ovh_cloud_managed_analytics_kafka_acl.acl.permission
}

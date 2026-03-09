data "ovh_cloud_managed_analytics_kafka_schemaregistryacl" "schema_registry_acl" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  id            = "ZZZ"
}

output "acl_permission" {
  value = data.ovh_cloud_managed_analytics_kafka_schemaregistryacl.schema_registry_acl.permission
}

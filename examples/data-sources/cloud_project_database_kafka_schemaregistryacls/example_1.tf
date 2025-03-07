data "ovh_cloud_project_database_kafka_schemaregistryacls" "schema_registry_acls" {
  service_name  = "XXX"
  cluster_id    = "YYY"
}

output "acl_ids" {
  value = data.ovh_cloud_project_database_kafka_schemaregistryacls.schema_registry_acls.acl_ids
}

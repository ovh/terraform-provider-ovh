data "ovh_cloud_project_database_kafka_acl" "acl" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  id            = "ZZZ"
}

output "acl_permission" {
  value = data.ovh_cloud_project_database_kafka_acl.acl.permission
}

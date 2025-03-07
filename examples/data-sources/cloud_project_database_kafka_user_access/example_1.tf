data "ovh_cloud_project_database_kafka_user_access" "access" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  user_id       = "ZZZ"
}

output "access_cert" {
  value = data.ovh_cloud_project_database_kafka_user_access.access.cert
}

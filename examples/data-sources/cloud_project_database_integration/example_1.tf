data "ovh_cloud_project_database_integration" "integration" {
  service_name  = "XXX"
  engine        = "YYY"
  cluster_id    = "ZZZ"
  id            = "UUU"
}

output "integration_type" {
  value = data.ovh_cloud_project_database_integration.integration.type
}

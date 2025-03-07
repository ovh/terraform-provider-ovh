data "ovh_cloud_project_database_integrations" "integrations" {
  service_name  = "XXX"
  engine        = "YYY"
  cluster_id    = "ZZZ"
}

output "integration_ids" {
  value = data.ovh_cloud_project_database_integrations.integrations.integration_ids
}

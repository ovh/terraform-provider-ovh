data "ovh_cloud_managed_analytics_integrations" "integrations" {
  service_name  = "XXX"
  engine        = "YYY"
  cluster_id    = "ZZZ"
}

output "integration_ids" {
  value = data.ovh_cloud_managed_analytics_integrations.integrations.integration_ids
}

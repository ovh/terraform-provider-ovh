data "ovh_cloud_managed_analytics_integration" "integration" {
  service_name  = "XXX"
  engine        = "YYY"
  cluster_id    = "ZZZ"
  id            = "UUU"
}

output "integration_type" {
  value = data.ovh_cloud_managed_analytics_integration.integration.type
}

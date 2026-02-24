data "ovh_cloud_managed_analytics_opensearch_pattern" "pattern" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  id            = "ZZZ"
}

output "pattern_pattern" {
  value = data.ovh_cloud_managed_analytics_opensearch_pattern.pattern.pattern
}

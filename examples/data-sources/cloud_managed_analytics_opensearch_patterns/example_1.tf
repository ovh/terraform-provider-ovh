data "ovh_cloud_managed_analytics_opensearch_patterns" "patterns" {
  service_name  = "XXX"
  cluster_id    = "YYY"
}

output "pattern_ids" {
  value = data.ovh_cloud_managed_analytics_opensearch_patterns.patterns.pattern_ids
}

data "ovh_cloud_project_database_opensearch_patterns" "patterns" {
  service_name  = "XXX"
  cluster_id    = "YYY"
}

output "pattern_ids" {
  value = data.ovh_cloud_project_database_opensearch_patterns.patterns.pattern_ids
}

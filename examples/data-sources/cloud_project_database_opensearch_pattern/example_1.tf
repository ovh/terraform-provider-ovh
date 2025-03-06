data "ovh_cloud_project_database_opensearch_pattern" "pattern" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  id            = "ZZZ"
}

output "pattern_pattern" {
  value = data.ovh_cloud_project_database_opensearch_pattern.pattern.pattern
}

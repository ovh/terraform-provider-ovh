data "ovh_cloud_project_database" "opensearch" {
  service_name  = "XXX"
  engine        = "opensearch"
  id            = "ZZZ"
}

resource "ovh_cloud_project_database_opensearch_pattern" "pattern" {
  service_name    = data.ovh_cloud_project_database.opensearch.service_name
  cluster_id      = data.ovh_cloud_project_database.opensearch.id
  max_index_count = 2
  pattern         = "logs_*"
}

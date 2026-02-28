data "ovh_cloud_managed_analytics" "db_postgresql" {
  service_name  = "XXXX"
  engine        = "postgresql"
  id            = "ZZZZ"
}

data "ovh_cloud_managed_analytics" "db_opensearch" {
  service_name  = "XXXX"
  engine        = "opensearch"
  id            = "ZZZZ"
}

resource "ovh_cloud_managed_analytics_integration" "integration" {
  service_name            = data.ovh_cloud_managed_analytics.db_postgresql.service_name
  engine                  = data.ovh_cloud_managed_analytics.db_postgresql.engine
  cluster_id              = data.ovh_cloud_managed_analytics.db_postgresql.id
  source_service_id       = data.ovh_cloud_managed_analytics.db_postgresql.id
  destination_service_id  = data.ovh_cloud_managed_analytics.db_opensearch.id
  type                    = "opensearchLogs"
}

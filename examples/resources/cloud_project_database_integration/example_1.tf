data "ovh_cloud_project_database" "db_postgresql" {
  service_name  = "XXXX"
  engine        = "postgresql"
  id            = "ZZZZ"
}

data "ovh_cloud_project_database" "db_opensearch" {
  service_name  = "XXXX"
  engine        = "opensearch"
  id            = "ZZZZ"
}

resource "ovh_cloud_project_database_integration" "integration" {
  service_name            = data.ovh_cloud_project_database.db_postgresql.service_name
  engine                  = data.ovh_cloud_project_database.db_postgresql.engine
  cluster_id              = data.ovh_cloud_project_database.db_postgresql.id
  source_service_id       = data.ovh_cloud_project_database.db_postgresql.id
  destination_service_id  = data.ovh_cloud_project_database.db_opensearch.id
  type                    = "opensearchLogs"
}

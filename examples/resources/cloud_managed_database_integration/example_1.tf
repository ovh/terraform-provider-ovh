resource "ovh_cloud_managed_database_integration" "integration" {
  service_name           = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  engine                 = "postgresql"
  service_id             = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  source_service_id      = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  destination_service_id = "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz"
  type                   = "GRAFANA_DATASOURCE"
}

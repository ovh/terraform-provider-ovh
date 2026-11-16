resource "ovh_cloud_managed_analytics_integration" "integration" {
  service_name           = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  engine                 = "kafka"
  service_id             = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  source_service_id      = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  destination_service_id = "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz"
  type                   = "KAFKA_LOGS"
}

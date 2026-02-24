data "ovh_cloud_managed_analytics" "kafka" {
  service_name  = "XXX"
  engine        = "kafka"
  id            = "ZZZ"
}

resource "ovh_cloud_managed_analytics_kafka_schemaregistryacl" "schema_registry_acl" {
  service_name    = data.ovh_cloud_managed_analytics.kafka.service_name
  cluster_id      = data.ovh_cloud_managed_analytics.kafka.id
  permission      = "schema_registry_read"
  resource        = "Subject:myResource"
  username        = "johndoe"
}

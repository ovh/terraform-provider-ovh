data "ovh_cloud_managed_analytics" "kafka" {
  service_name  = "XXX"
  engine        = "kafka"
  id            = "ZZZ"
}

resource "ovh_cloud_managed_analytics_kafka_topic" "topic" {
  service_name        = data.ovh_cloud_managed_analytics.kafka.service_name
  cluster_id          = data.ovh_cloud_managed_analytics.kafka.id
  name                = "mytopic"
  min_insync_replicas = 1
  partitions          = 3
  replication         = 2
  retention_bytes     = 4
  retention_hours     = 5
}

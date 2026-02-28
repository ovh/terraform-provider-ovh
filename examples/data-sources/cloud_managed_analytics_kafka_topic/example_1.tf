data "ovh_cloud_managed_analytics_kafka_topic" "topic" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  id            = "ZZZ"
}

output "topic_name" {
  value = data.ovh_cloud_managed_analytics_kafka_topic.topic.name
}

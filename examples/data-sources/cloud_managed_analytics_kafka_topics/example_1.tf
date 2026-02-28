data "ovh_cloud_managed_analytics_kafka_topics" "topics" {
  service_name  = "XXX"
  cluster_id    = "YYY"
}

output "topic_ids" {
  value = data.ovh_cloud_managed_analytics_kafka_topics.topics.topic_ids
}

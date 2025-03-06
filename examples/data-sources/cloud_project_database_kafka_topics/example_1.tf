data "ovh_cloud_project_database_kafka_topics" "topics" {
  service_name  = "XXX"
  cluster_id    = "YYY"
}

output "topic_ids" {
  value = data.ovh_cloud_project_database_kafka_topics.topics.topic_ids
}

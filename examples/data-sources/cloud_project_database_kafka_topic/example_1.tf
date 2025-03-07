data "ovh_cloud_project_database_kafka_topic" "topic" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  id            = "ZZZ"
}

output "topic_name" {
  value = data.ovh_cloud_project_database_kafka_topic.topic.name
}

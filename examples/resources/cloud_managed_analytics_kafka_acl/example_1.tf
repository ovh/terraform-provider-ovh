data "ovh_cloud_managed_analytics" "kafka" {
  service_name  = "XXX"
  engine        = "kafka"
  id            = "ZZZ"
}

resource "ovh_cloud_managed_analytics_kafka_acl" "acl" {
  service_name    = data.ovh_cloud_managed_analytics.kafka.service_name
  cluster_id      = data.ovh_cloud_managed_analytics.kafka.id
  permission      = "read"
  topic           = "mytopic"
  username        = "johndoe"
}

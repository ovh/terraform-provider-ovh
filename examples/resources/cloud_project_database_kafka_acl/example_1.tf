data "ovh_cloud_project_database" "kafka" {
  service_name  = "XXX"
  engine        = "kafka"
  id            = "ZZZ"
}

resource "ovh_cloud_project_database_kafka_acl" "acl" {
  service_name    = data.ovh_cloud_project_database.kafka.service_name
  cluster_id      = data.ovh_cloud_project_database.kafka.id
  permission      = "read"
  topic           = "mytopic"
  username        = "johndoe"
}

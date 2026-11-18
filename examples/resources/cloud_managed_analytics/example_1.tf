resource "ovh_cloud_managed_analytics" "kafka" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  engine       = "kafka"
  description  = "my-kafka"
  region       = "eu-west-gra"
  version      = "3.8"
  plan         = "business"
  flavor       = "b3-8"
  nodes_count  = 3
}

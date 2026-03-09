resource "ovh_cloud_managed_analytics" "kafka" {
  service_name          = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description           = "my-first-kafka"
  engine                = "kafka"
  version               = "3.8"
  plan                  = "business"
  kafka_rest_api        = true
  kafka_schema_registry = true
  nodes {
    region  = "GRA"
  }
  nodes {
    region  = "GRA"
  }
  nodes {
    region  = "GRA"
  }
  flavor                = "db1-15"
}

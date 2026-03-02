resource "ovh_cloud_managed_analytics" "kafkadb" {
  service_name          = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description           = "my-first-kafka"
  engine                = "kafka"
  version               = "3.8"
  flavor                = "db1-4"
  plan                  = "business"
  kafka_rest_api        = true
  kafka_schema_registry = true
  nodes {
    region  = "DE"
  }
  nodes {
    region  = "DE"
  }
  nodes {
    region  = "DE"
  }
}

resource "ovh_cloud_managed_analytics" "opensearchdb" {
  service_name            = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description             = "my-first-opensearch"
  engine                  = "opensearch"
  version                 = "1"
  plan                    = "essential"
  opensearch_acls_enabled = true
  nodes {
    region = "UK"
  }
  flavor                  = "db1-4"
}

resource "ovh_cloud_managed_analytics" "grafana" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-grafana"
  engine        = "grafana"
  version       = "9.1"
  plan          = "essential"
  nodes {
    region =  "GRA"
  }
  flavor        = "db1-4"
}

resource "ovh_cloud_project_database" "cassandradb" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-cassandra"
  engine        = "cassandra"
  version       = "4.0"
  plan          = "essential"
  nodes {
    region  = "BHS"
  }
  nodes {
    region  = "BHS"
  }
  nodes {
    region  = "BHS"
  }
  flavor        = "db1-4"
}

resource "ovh_cloud_project_database" "kafkadb" {
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

resource "ovh_cloud_project_database" "m3db" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-m3db"
  engine        = "m3db"
  version       = "1.2"
  plan          = "essential"
  nodes {
    region  = "BHS"
  }
  flavor        = "db1-7"
}

resource "ovh_cloud_project_database" "mongodb" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-mongodb"
  engine        = "mongodb"
  version       = "5.0"
  plan          = "discovery"
  nodes {
    region =  "GRA"
  }
  flavor        = "db1-2"
}

resource "ovh_cloud_project_database" "mysqldb" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-mysql"
  engine        = "mysql"
  version       = "8"
  plan          = "essential"
  nodes {
    region  = "SBG"
  }
  flavor        = "db1-4"
  advanced_configuration = {
    "mysql.sql_mode": "ANSI,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION,NO_ZERO_DATE,NO_ZERO_IN_DATE,STRICT_ALL_TABLES",
    "mysql.sql_require_primary_key": "true"
  }
}

resource "ovh_cloud_project_database" "opensearchdb" {
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

resource "ovh_cloud_project_database" "pgsqldb" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-postgresql"
  engine        = "postgresql"
  version       = "14"
  plan          = "essential"
  nodes {
    region  = "WAW"
  }
  flavor        = "db1-4"
  ip_restrictions {
    description = "ip 1"
    ip = "178.97.6.0/24"
  }
  ip_restrictions {
    description = "ip 2"
    ip = "178.97.7.0/24"
  }
}

resource "ovh_cloud_project_database" "redisdb" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-redis"
  engine        = "redis"
  version       = "6.2"
  plan          = "essential"
  nodes {
    region  = "BHS"
  }
  flavor        = "db1-4"
}

resource "ovh_cloud_project_database" "grafana" {
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

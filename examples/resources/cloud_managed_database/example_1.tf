resource "ovh_cloud_managed_database" "pgsqldb" {
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

resource "ovh_cloud_managed_database" "mysqldb" {
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

resource "ovh_cloud_managed_database" "mongodb" {
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

resource "ovh_cloud_managed_database" "valkeydb" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-valkey"
  engine        = "valkey"
  version       = "8.0"
  plan          = "essential"
  nodes {
    region  = "BHS"
  }
  flavor        = "db1-4"
}

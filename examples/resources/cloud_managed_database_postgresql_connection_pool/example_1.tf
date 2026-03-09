resource "ovh_cloud_managed_database" "db" {
  service_name  = "XXXX"
  engine        = "postgresql"
  description  = "test-postgresql-cluster"
  version      = "15"
  plan         = "essential"
  nodes {
    region     = "GRA"
  }
  flavor = "db1-4"
}

resource "ovh_cloud_managed_database_database" "database" {
  service_name  = ovh_cloud_managed_database.db.service_name
  engine        = ovh_cloud_managed_database.db.engine
  cluster_id    = ovh_cloud_managed_database.db.id
  name          = "mydatabase"
}

resource "ovh_cloud_managed_database_postgresql_user" "user" {
  service_name = ovh_cloud_managed_database.db.service_name
  cluster_id   = ovh_cloud_managed_database.db.id
  name          = "johndoe"
  roles         = ["replication"]
}

resource "ovh_cloud_managed_database_postgresql_connection_pool" "test_pool" {
  service_name = ovh_cloud_managed_database.db.service_name
  cluster_id   = ovh_cloud_managed_database.db.id
  database_id = ovh_cloud_managed_database_database.database.id
  name = "test_connection_pool"
  user_id = ovh_cloud_managed_database_postgresql_user.user.id
  mode = "session"
  size = 13
}

output "test_pool" {
  value = {
    service_name: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.service_name
    cluster_id: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.cluster_id
    name: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.name
    database_id: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.database_id
    mode: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.mode
    size: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.size
    port: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.port
    ssl_mode: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.ssl_mode
    uri: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.uri
    user_id: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.user_id
  }
}

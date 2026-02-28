data "ovh_cloud_managed_database_postgresql_connection_pool" "test_pool" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ"
}

output "test_pool" {
  value = {
    service_name: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.service_name
    cluster_id: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.cluster_id
    name: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.name
    database_id: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.database_id
    mode: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.mode
    size: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.size
    port: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.port
    ssl_mode: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.ssl_mode
    uri: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.uri
    user_id: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.user_id
  }
}

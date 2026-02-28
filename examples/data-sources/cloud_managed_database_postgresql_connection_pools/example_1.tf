data "ovh_cloud_managed_database_postgresql_connection_pools" "test_pools" {
  service_name  = "XXX"
  cluster_id    = "YYY"
}

output "connection_pool_ids" {
  value = data.ovh_cloud_managed_database_postgresql_connection_pools.test_pools.connection_pool_ids
}

data "ovh_cloud_managed_database_databases" "databases" {
  service_name  = "XXXX"
  engine        = "YYYY"
  cluster_id    = "ZZZ"
}

output "database_ids" {
  value = data.ovh_cloud_managed_database_databases.databases.database_ids
}

data "ovh_cloud_managed_databases" "dbs" {
  service_name  = "XXXXXX"
  engine        = "YYYY"
}

output "cluster_ids" {
  value = data.ovh_cloud_managed_databases.dbs.cluster_ids
}

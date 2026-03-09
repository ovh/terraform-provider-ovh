data "ovh_cloud_managed_database_database" "database" {
  service_name  = "XXX"
  engine        = "YYY"
  cluster_id    = "ZZZ"
  name          = "UUU"
}

output "database_name" {
  value = data.ovh_cloud_managed_database_database.database.name
}

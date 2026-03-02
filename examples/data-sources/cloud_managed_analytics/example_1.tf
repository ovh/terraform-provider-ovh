data "ovh_cloud_managed_analytics" "db" {
  service_name  = "XXXXXX"
  engine        = "YYYY"
  id            = "ZZZZ"
}

output "cluster_id" {
  value = data.ovh_cloud_managed_analytics.db.id
}

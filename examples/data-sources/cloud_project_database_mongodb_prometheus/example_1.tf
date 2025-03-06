data "ovh_cloud_project_database_mongodb_prometheus" "prometheus" {
  service_name  = "XXX"
  cluster_id    = "ZZZ"
}

output "name" {
  value = data.ovh_cloud_project_database_mongodb_prometheus.prometheus.username
}

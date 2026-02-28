data "ovh_cloud_managed_analytics_prometheus" "prometheus" {
  service_name  = "XXX"
  engine        = "YYY"
  cluster_id    = "ZZZ"
}

output "name" {
  value = data.ovh_cloud_managed_analytics_prometheus.prometheus.username
}

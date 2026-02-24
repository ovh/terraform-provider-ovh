data "ovh_cloud_managed_analytics_m3db_namespace" "m3db_namespace" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ"
}

output "m3dbnamespace_type" {
  value = data.ovh_cloud_managed_analytics_m3db_namespace.m3db_namespace.type
}

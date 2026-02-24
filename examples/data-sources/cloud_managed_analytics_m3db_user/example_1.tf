data "ovh_cloud_managed_analytics_m3db_user" "m3db_user" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ"
}

output "m3db_user_group" {
  value = data.ovh_cloud_managed_analytics_m3db_user.m3db_user.group
}

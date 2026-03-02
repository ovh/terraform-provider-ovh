data "ovh_cloud_managed_analytics_user" "user" {
  service_name  = "XXX"
  engine        = "YYY"
  cluster_id    = "ZZZ"
  name          = "UUU"
}

output "user_name" {
  value = data.ovh_cloud_managed_analytics_user.user.name
}

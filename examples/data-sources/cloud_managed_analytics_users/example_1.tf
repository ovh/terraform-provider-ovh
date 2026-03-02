data "ovh_cloud_managed_analytics_users" "users" {
  service_name  = "XXXX"
  engine        = "YYYY"
  cluster_id    = "ZZZ"
}

output "user_ids" {
  value = data.ovh_cloud_managed_analytics_users.users.user_ids
}

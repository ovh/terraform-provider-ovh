data "ovh_cloud_managed_database_redis_user" "redis_user" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ"
}

output "redis_user_commands" {
  value = data.ovh_cloud_managed_database_redis_user.redis_user.commands
}

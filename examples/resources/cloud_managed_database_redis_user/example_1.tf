data "ovh_cloud_managed_database" "redis" {
  service_name  = "XXXX"
  engine        = "redis"
  id            = "ZZZZ"
}

resource "ovh_cloud_managed_database_redis_user" "user" {
  service_name  = data.ovh_cloud_managed_database.redis.service_name
  cluster_id    = data.ovh_cloud_managed_database.redis.id
  categories    = ["+@set", "+@sortedset"]
  channels      = ["*"]
  commands      = ["+get", "-set"]
  keys          = ["data", "properties"]
  name          = "johndoe"
}

output "user_password" {
  value     = ovh_cloud_managed_database_redis_user.user.password
  sensitive = true
}

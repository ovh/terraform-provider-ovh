data "ovh_cloud_managed_analytics" "db" {
  service_name  = "XXXX"
  engine        = "YYYY"
  id            = "ZZZZ"
}

resource "ovh_cloud_managed_analytics_user" "user" {
  service_name  = data.ovh_cloud_managed_analytics.db.service_name
  engine        = data.ovh_cloud_managed_analytics.db.engine
  cluster_id    = data.ovh_cloud_managed_analytics.db.id
  name          = "johndoe"
}

output "user_password" {
  value     = ovh_cloud_managed_analytics_user.user.password
  sensitive = true
}

data "ovh_cloud_managed_analytics" "db" {
  service_name = "XXXX"
  engine       = "YYYY"
  id           = "ZZZZ"
}

resource "ovh_cloud_managed_analytics_ip_restriction" "ip_restriction" {
  service_name = data.ovh_cloud_managed_analytics.db.service_name
  engine       = data.ovh_cloud_managed_analytics.db.engine
  cluster_id   = data.ovh_cloud_managed_analytics.db.id
  ip           = "178.97.6.0/24"
}

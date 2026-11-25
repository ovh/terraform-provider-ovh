resource "ovh_cloud_managed_analytics_prometheus" "prometheus" {
  service_name   = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  engine         = "kafka"
  service_id     = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  password_reset = "2026-11-12T10:00:00Z"
}

output "prometheus_password" {
  value     = ovh_cloud_managed_analytics_prometheus.prometheus.password
  sensitive = true
}

data "ovh_cloud_managed_analytics" "db" {
  service_name  = "XXXX"
  engine        = "YYYY"
  id            = "ZZZZ"
}

# Set password_reset to be based on the update of another variable to reset the password
resource "ovh_cloud_managed_analytics_prometheus" "prometheus_datetime" {
  service_name    = data.ovh_cloud_managed_analytics.db.service_name
  engine          = data.ovh_cloud_managed_analytics.db.engine
  cluster_id      = data.ovh_cloud_managed_analytics.db.id
  password_reset  = "2024-01-02T11:00:00Z"
}

variable "something" {
  type = string
}

resource "ovh_cloud_managed_analytics_prometheus" "prometheus_md5" {
  service_name    = data.ovh_cloud_managed_analytics.db.service_name
  engine          = data.ovh_cloud_managed_analytics.db.engine
  cluster_id      = data.ovh_cloud_managed_analytics.db.id
  password_reset  = md5(var.something)
}

resource "ovh_cloud_managed_analytics_prometheus" "prometheus" {
  service_name    = data.ovh_cloud_managed_analytics.db.service_name
  engine          = data.ovh_cloud_managed_analytics.db.engine
  cluster_id      = data.ovh_cloud_managed_analytics.db.id
  password_reset  = "reset1"
}

output "prom_password" {
  value     = ovh_cloud_managed_analytics_prometheus.prometheus.password
  sensitive = true
}

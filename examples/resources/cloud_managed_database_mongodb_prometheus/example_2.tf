data "ovh_cloud_managed_database" "mongodb" {
  service_name  = "XXXX"
  engine        = "mongodb"
  id            = "ZZZZ"
}

# Change password_reset with the datetime each time you want to reset the password to trigger an update
resource "ovh_cloud_managed_database_mongodb_prometheus" "prometheus_datetime" {
  service_name    = data.ovh_cloud_managed_database.mongodb.service_name
  cluster_id      = data.ovh_cloud_managed_database.mongodb.id
  password_reset  = "2024-01-02T11:00:00Z"
}

# Set password_reset to be based on the update of another variable to reset the password
resource "ovh_cloud_managed_database_mongodb_prometheus" "prometheus_md5" {
  service_name    = data.ovh_cloud_managed_database.mongodb.service_name
  cluster_id      = data.ovh_cloud_managed_database.mongodb.id
  password_reset  = md5(var.something)
}

# Change password_reset each time you want to reset the password to trigger an update
resource "ovh_cloud_managed_database_mongodb_prometheus" "prometheus" {
  service_name    = data.ovh_cloud_managed_database.mongodb.service_name
  cluster_id      = data.ovh_cloud_managed_database.mongodb.id
  password_reset  = "reset1"
}

output "prom_password" {
  value     = ovh_cloud_managed_database_mongodb_prometheus.prometheus.password
  sensitive = true
}

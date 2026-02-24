data "ovh_cloud_managed_database" "db" {
  service_name  = "XXXX"
  engine        = "YYYY"
  id            = "ZZZZ"
}

resource "ovh_cloud_managed_database_prometheus" "prometheus" {
  service_name  = data.ovh_cloud_managed_database.db.service_name
  engine        = data.ovh_cloud_managed_database.db.engine
  cluster_id    = data.ovh_cloud_managed_database.db.id
}

output "prom_password" {
  value     = ovh_cloud_managed_database_prometheus.prometheus.password
  sensitive = true
}

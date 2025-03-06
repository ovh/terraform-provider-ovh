data "ovh_cloud_project_database" "mongodb" {
  service_name  = "XXX"
  engine        = "mongodb"
  id            = "ZZZ"
}

resource "ovh_cloud_project_database_mongodb_prometheus" "prometheus" {
  service_name  = data.ovh_cloud_project_database.mongodb.service_name
  cluster_id    = data.ovh_cloud_project_database.mongodb.id
}

output "prom_password" {
  value     = ovh_cloud_project_database_mongodb_prometheus.prometheus.password
  sensitive = true
}

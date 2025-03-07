data "ovh_cloud_project_database" "db" {
  service_name = "XXXX"
  engine       = "YYYY"
  id           = "ZZZZ"
}

resource "ovh_cloud_project_database_ip_restriction" "ip_restriction" {
  service_name = data.ovh_cloud_project_database.db.service_name
  engine       = data.ovh_cloud_project_database.db.engine
  cluster_id   = data.ovh_cloud_project_database.db.id
  ip           = "178.97.6.0/24"
}

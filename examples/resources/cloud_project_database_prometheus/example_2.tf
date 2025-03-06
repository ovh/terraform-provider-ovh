data "ovh_cloud_project_database" "db" {
  service_name  = "XXXX"
  engine        = "YYYY"
  id            = "ZZZZ"
}

# Set password_reset to be based on the update of another variable to reset the password
resource "ovh_cloud_project_database_prometheus" "prometheusDatetime" {
  service_name    = data.ovh_cloud_project_database.db.service_name
  engine          = data.ovh_cloud_project_database.db.engine
  cluster_id      = data.ovh_cloud_project_database.db.id
  password_reset  = "2024-01-02T11:00:00Z"
}

variable "something" {
  type = string
}

resource "ovh_cloud_project_database_prometheus" "prometheusMd5" {
  service_name    = data.ovh_cloud_project_database.db.service_name
  engine          = data.ovh_cloud_project_database.db.engine
  cluster_id      = data.ovh_cloud_project_database.db.id
  password_reset  = md5(var.something)
}

resource "ovh_cloud_project_database_prometheus" "prometheus" {
  service_name    = data.ovh_cloud_project_database.db.service_name
  engine          = data.ovh_cloud_project_database.db.engine
  cluster_id      = data.ovh_cloud_project_database.db.id
  password_reset  = "reset1"
}

output "prom_password" {
  value     = ovh_cloud_project_database_prometheus.prometheus.password
  sensitive = true
}

data "ovh_cloud_project_database" "db" {
  service_name  = "XXXXXX"
  engine        = "YYYY"
  id            = "ZZZZ"
}

output "cluster_id" {
  value = data.ovh_cloud_project_database.db.id
}

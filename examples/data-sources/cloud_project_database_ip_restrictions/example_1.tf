data "ovh_cloud_project_database_ip_restrictions" "ip_restrictions" {
  service_name  = "XXXXXX"
  engine        = "YYYY"
  cluster_id    = "ZZZZ"
}

output "ips" {
  value = data.ovh_cloud_project_database_ip_restrictions.ip_restrictions.ips
}

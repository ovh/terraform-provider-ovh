data "ovh_cloud_project_database_m3db_namespaces" "namespaces" {
  service_name  = "XXX"
  cluster_id    = "YYY"
}

output "namespace_ids" {
  value = data.ovh_cloud_project_database_m3db_namespaces.namespaces.namespace_ids
}

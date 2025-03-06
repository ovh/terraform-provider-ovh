data "ovh_cloud_project_database_m3db_namespace" "m3db_namespace" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ"
}

output "m3dbnamespace_type" {
  value = data.ovh_cloud_project_database_m3db_namespace.m3db_namespace.type
}

data "ovh_cloud_project_database_capabilities" "capabilities" {
  service_name  = "XXX"
}

output "capabilities_engine_name" {
  value = tolist(data.ovh_cloud_project_database_capabilities.capabilities[*].engines)[0]
}

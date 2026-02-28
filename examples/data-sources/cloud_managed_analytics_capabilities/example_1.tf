data "ovh_cloud_managed_analytics_capabilities" "capabilities" {
  service_name  = "XXX"
}

output "capabilities_engine_name" {
  value = tolist(data.ovh_cloud_managed_analytics_capabilities.capabilities[*].engines)[0]
}

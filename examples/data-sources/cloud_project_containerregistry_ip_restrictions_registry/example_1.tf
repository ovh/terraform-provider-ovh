data data "ovh_cloud_project_containerregistry_ip_restrictions_registry" "my_iprestrictions_data" {
  service_name = "XXXXXX"
  registry_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}

output "my_ip_restrictions" {
  value = data.ovh_cloud_project_containerregistry_ip_restrictions_registry.my_iprestrictions_data.ip_restrictions
}

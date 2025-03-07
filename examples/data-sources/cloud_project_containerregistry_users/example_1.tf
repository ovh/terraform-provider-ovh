data "ovh_cloud_project_containerregistry" "my_registry" {
  service_name = "XXXXXX"
  registry_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}

data "ovh_cloud_project_containerregistry_users" "users" {
  service_name = ovh_cloud_project_containerregistry.my_registry.service_name
  registry_id  = ovh_cloud_project_containerregistry.my_registry.id
}

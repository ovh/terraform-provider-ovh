data "ovh_cloud_project_containerregistry" "registry" {
  service_name = "XXXXXX"
  registry_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}

resource "ovh_cloud_project_containerregistry_user" "user" {
  service_name = ovh_cloud_project_containerregistry.registry.service_name
  registry_id  = ovh_cloud_project_containerregistry.registry.id
  email        = "foo@bar.com"
  login        = "foobar"
}

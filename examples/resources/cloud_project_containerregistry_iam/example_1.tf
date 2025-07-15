resource "ovh_cloud_project_containerregistry_iam" "registry_iam" {
  service_name = "XXXXXX"
  registry_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"

  #optional field
  delete_users = false
}

output "iam_enabled" {
  value     = ovh_cloud_project_containerregistry_iam.registry_iam.iam_enabled
  sensitive = true
}

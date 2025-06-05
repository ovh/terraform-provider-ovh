resource "ovh_cloud_project_containerregistry_iam" "my_iam" {
  service_name = "XXXXXX"
  registry_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"

  #optional field
  delete_users = false
}

output "iam-enabled" {
  value     = ovh_cloud_project_containerregistry_iam.my_iam.iam_enabled
  sensitive = true
}

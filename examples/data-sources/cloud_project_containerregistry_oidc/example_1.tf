data "ovh_cloud_project_containerregistry_oidc" "my_oidc" {
  service_name = "XXXXXX"
  registry_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}

output "oidc-client-id" {
  value = data.ovh_cloud_project_containerregistry_oidc.my_oidc.oidc_client_id
}

resource "ovh_cloud_project_containerregistry_oidc" "my_oidc" {
  service_name = "XXXXXX"
  registry_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"

  #required field
  oidc_name          = "my-oidc-provider"
  oidc_endpoint      = "https://xxxx.yyy.com"
  oidc_client_id     = "xxx"
  oidc_client_secret = "xxx"
  oidc_scope         = "openid,profile,email,offline_access"

  #optional field
  oidc_groups_claim = "groups"
  oidc_admin_group  = "harbor-admin"
  oidc_verify_cert  = true
  oidc_auto_onboard = true
  oidc_user_claim   = "preferred_username"
  delete_users      = false
}

output "oidc_client_secret" {
  value = ovh_cloud_project_containerregistry_oidc.my_oidc.oidc_client_secret
  sensitive = true
}

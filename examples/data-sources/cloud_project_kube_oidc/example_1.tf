data "ovh_cloud_project_kube_oidc" "oidc" {
  service_name = "XXXXXX"
  kube_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}

output "oidc-val" {
  value = data.ovh_cloud_project_kube_oidc.oidc.client_id
}

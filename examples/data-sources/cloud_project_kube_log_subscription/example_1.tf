data "ovh_cloud_project_kube_log_subscription" "sub" {
  service_name    = "XXXXXX"
  kube_id         = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  subscription_id = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
}

output "resource-name" {
  value = data.ovh_cloud_project_kube_log_subscription.sub.resource.0.name
}

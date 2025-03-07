data "ovh_cloud_project_kube_iprestrictions" "ip_restrictions" {
  service_name = "XXXXXX"
  kube_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}

output "ips" {
  value = data.ovh_cloud_project_kube_iprestrictions.ip_restrictions.ips
}

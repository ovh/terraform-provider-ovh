data "ovh_cloud_managed_analytics_certificates" "certificates" {
  service_name  = "XXX"
  engine        = "YYY"
  cluster_id    = "ZZZ"
}

output "certificates_ca" {
  value = data.ovh_cloud_managed_analytics_certificates.certificates.ca
}

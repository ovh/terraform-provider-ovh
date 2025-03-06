data "ovh_cloud_project_database_certificates" "certificates" {
  service_name  = "XXX"
  engine        = "YYY"
  cluster_id    = "ZZZ"
}

output "certificates_ca" {
  value = data.ovh_cloud_project_database_certificates.certificates.ca
}

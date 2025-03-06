data "ovh_cloud_project_vrack" "vrack" {
  service_name = "XXXXXX"
}

output "vrack" {
  value = data.ovh_cloud_project_vrack.vrack
}

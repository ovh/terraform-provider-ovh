data "ovh_cloud_managed_kubernetes_plans" "plans" {
  service_name = "<Public cloud project id>"
  region       = "EU-WEST-PAR"
}

output "available_plans" {
  value = [for p in data.ovh_cloud_managed_kubernetes_plans.plans.plans : p.name if p.status == "AVAILABLE"]
}

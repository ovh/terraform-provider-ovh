data "ovh_cloud_managed_kubernetes_regions" "regions" {
  service_name = "<Public cloud project id>"
}

# Regions offering the STANDARD plan.
output "standard_regions" {
  value = [for r in data.ovh_cloud_managed_kubernetes_regions.regions.regions : r.name if contains(r.available_plans, "STANDARD")]
}

data "ovh_cloud_managed_kubernetes_capabilities_versions" "versions" {
  service_name          = "<Public cloud project id>"
  managed_kubernetes_id = "<cluster id>"
}

output "current_version" {
  value = data.ovh_cloud_managed_kubernetes_capabilities_versions.versions.current_version
}

# Each list holds at most one version: the highest one available.
output "latest_patch_version" {
  value = try(data.ovh_cloud_managed_kubernetes_capabilities_versions.versions.available_updates[0].version, data.ovh_cloud_managed_kubernetes_capabilities_versions.versions.current_version)
}

output "next_minor_version" {
  value = try(data.ovh_cloud_managed_kubernetes_capabilities_versions.versions.available_upgrades[0].version, null)
}

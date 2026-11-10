---
subcategory: "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_managed_kubernetes_capabilities_versions (Data Source)

Use this data source to get the Kubernetes version deployed on a Managed Kubernetes cluster and the versions it can move to.

## Example Usage

```terraform
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
```

## Argument Reference

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `managed_kubernetes_id` - (Required) ID of the [`ovh_cloud_managed_kubernetes`](../resources/cloud_managed_kubernetes.md) cluster.

## Attributes Reference

* `current_version` - Full Kubernetes version currently deployed on the cluster (for example `1.32.8-1`).
* `available_updates` - Patch versions of the current minor the cluster can update to (for example `1.32.8-1` to `1.32.8-2`):
  * `version` - Full Kubernetes version (for example `1.32.8-2`).
  * `status` - Support status of the version (`SUPPORTED`, `END_OF_SUPPORT`).
* `available_upgrades` - Versions of the next minor the cluster can upgrade to (for example `1.32.8-1` to `1.33.1-2`). Multi-minor jumps are never listed:
  * `version` - Full Kubernetes version (for example `1.33.1-2`).
  * `status` - Support status of the version (`SUPPORTED`, `END_OF_SUPPORT`).
